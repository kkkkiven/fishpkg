package kafka

import (
	"hash/fnv"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/logs"

	"github.com/Shopify/sarama"

	cluster "github.com/bsm/sarama-cluster"

	"github.com/pkg/errors"
)

type Option struct {
	Brokers []string
	Group   string
	Topics  []string
}

type Consumer struct {
	consumer *cluster.Consumer
	option   *Option
	chClose  chan struct{}
}

func NewConsumer(opt *Option) (*Consumer, error) {
	if opt == nil {
		return nil, errors.Errorf("kafka option is <nil>")
	}

	conf := cluster.NewConfig()
	conf.Net.TLS.Enable = false
	conf.Consumer.Return.Errors = false
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Group.Return.Notifications = false
	conf.Consumer.Offsets.CommitInterval = time.Duration(1) * time.Second
	// conf.Version = sarama.V2_2_0_0

	if err := conf.Validate(); err != nil {
		return nil, errors.WithMessage(err, "validate config failed")
	}

	consumer, err := cluster.NewConsumer(opt.Brokers, opt.Group, opt.Topics, conf)
	if err != nil {
		return nil, errors.Errorf("create kafka consumer error: %s, config: %v", err.Error(), conf)
	}

	c := &Consumer{}
	c.consumer = consumer
	c.option = new(Option)
	c.option.Brokers = append(c.option.Brokers, opt.Brokers...)
	c.option.Group = opt.Group
	c.option.Topics = append(c.option.Topics, opt.Topics...)
	c.chClose = make(chan struct{}, 0)

	return c, nil
}

func (this *Consumer) Start(msgHandle func(msg *sarama.ConsumerMessage) (bool, error)) {
	if this.consumer == nil {
		logs.Errorf("please new consumer first")
		return
	}

	logs.Infof("[%s] Start up kafka consumer ...", this.option.Group)

	for {
		select {
		case msg, ok := <-this.consumer.Messages():
			if ok {
				if ok, err := msgHandle(msg); err != nil {
					logs.Errorf("[%s] Handle error: %s", this.option.Group, err.Error())
				} else if ok {
					this.consumer.MarkOffset(msg, "")
				}
			}
			continue
		case err, ok := <-this.consumer.Errors():
			if ok {
				logs.Errorf("[%s] Kafka consumer error: %s", this.option.Group, err.Error())
			}
		case ntf, ok := <-this.consumer.Notifications(): // 群组中消费者数量变更通知
			if ok {
				logs.Debugf("[%s] Kafka consumer rebalance: %v", this.option.Group, ntf)
			}
		case <-this.chClose: // 服务停止
			logs.Infof("[%s] Stop kafka consumer ...", this.option.Group)
			if err := this.consumer.Close(); err != nil {
				logs.Errorf("[%s] Stop kafka consumer error: %s", this.option.Group, err.Error())
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (this *Consumer) Stop() {
	close(this.chClose)
}

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true
	// conf.Version = sarama.V2_2_0_0

	if err := conf.Validate(); err != nil {
		return nil, errors.WithMessage(err, "validate config failed")
	}

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		return nil, errors.WithMessage(err, "new sync producer failed")
	}

	return &Producer{producer: producer}, nil
}

func NewHashProducer(brokers []string) (*Producer, error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewCustomHashPartitioner(fnv.New32a)
	conf.Producer.Return.Successes = true
	// conf.Version = sarama.V2_2_0_0

	if err := conf.Validate(); err != nil {
		return nil, errors.WithMessage(err, "validate config failed")
	}

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		return nil, errors.WithMessage(err, "new sync producer failed")
	}

	return &Producer{producer: producer}, nil
}

func (this *Producer) SendMessage(topic string, message []byte, keys ...string) error {
	if this.producer == nil {
		return errors.Errorf("please new producer first")
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(message)
	if len(keys) != 0 {
		msg.Key = sarama.StringEncoder(keys[0])
	}

	if _, _, err := this.producer.SendMessage(msg); err != nil {
		return err
	}

	return nil
}

func (this *Producer) Reset(brokers []string) error {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewCustomHashPartitioner(fnv.New32a)
	conf.Producer.Return.Successes = true
	// conf.Version = sarama.V2_2_0_0

	if err := conf.Validate(); err != nil {
		return errors.WithMessage(err, "validate config failed")
	}

	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		return errors.WithMessage(err, "reset producer failed")
	}

	this.producer = producer
	return nil
}

func (this *Producer) Close() error {
	if this.producer == nil {
		return nil
	}

	return this.producer.Close()
}

type AsyncProducer struct {
	producer  sarama.AsyncProducer
	onSuccess func(*sarama.ProducerMessage)
	onError   func(*sarama.ProducerError)
	timeout   time.Duration
	chClose   chan struct{}
}

func NewAsyncProducer(brokers []string, timeout time.Duration, onSuccess func(*sarama.ProducerMessage), onError func(*sarama.ProducerError)) (*AsyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Timeout = timeout
	// conf.Version = sarama.V2_2_0_0

	producer, err := sarama.NewAsyncProducer(brokers, conf)
	if err != nil {
		return nil, errors.WithMessage(err, "new async producer failed")
	}

	chClose := make(chan struct{})

	errors := producer.Errors()
	success := producer.Successes()
	go func() {
		for {
			select {
			case err := <-errors:
				if onError != nil {
					onError(err)
				}
			case msg := <-success:
				if onSuccess != nil {
					onSuccess(msg)
				}
			case <-chClose:
				return
			}
		}
	}()

	return &AsyncProducer{producer: producer, onSuccess: onSuccess, onError: onError, timeout: timeout, chClose: chClose}, nil
}

func NewAsyncHashProducer(brokers []string, timeout time.Duration, onSuccess func(*sarama.ProducerMessage), onError func(*sarama.ProducerError)) (*AsyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Timeout = timeout
	conf.Producer.Partitioner = sarama.NewCustomHashPartitioner(fnv.New32a)
	// conf.Version = sarama.V2_2_0_0

	producer, err := sarama.NewAsyncProducer(brokers, conf)
	if err != nil {
		return nil, errors.WithMessage(err, "new async producer failed")
	}

	chClose := make(chan struct{})

	errors := producer.Errors()
	success := producer.Successes()
	go func() {
		for {
			select {
			case err := <-errors:
				if onError != nil {
					onError(err)
				}
			case msg := <-success:
				if onSuccess != nil {
					onSuccess(msg)
				}
			case <-chClose:
				return
			}
		}
	}()

	return &AsyncProducer{producer: producer, onSuccess: onSuccess, onError: onError, timeout: timeout, chClose: chClose}, nil
}

func (this *AsyncProducer) SendMessage(topic string, message []byte, keys ...string) error {
	if this.producer == nil {
		return errors.Errorf("please new producer first")
	}

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(message)
	if len(keys) != 0 {
		msg.Key = sarama.StringEncoder(keys[0])
	}

	this.producer.Input() <- msg
	return nil
}

func (this *AsyncProducer) Reset(brokers []string) error {
	if this == nil {
		return errors.New("please new producer first")
	}

	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Timeout = this.timeout
	// conf.Version = sarama.V2_2_0_0

	producer, err := sarama.NewAsyncProducer(brokers, conf)
	if err != nil {
		return errors.WithMessage(err, "new async producer failed")
	}

	if this.producer != nil {
		this.Close()
	}
	this.producer = producer
	this.chClose = make(chan struct{})

	errors := producer.Errors()
	success := producer.Successes()
	go func() {
		for {
			select {
			case err := <-errors:
				if this.onError != nil {
					this.onError(err)
				}
			case msg := <-success:
				if this.onSuccess != nil {
					this.onSuccess(msg)
				}
			case <-this.chClose:
				return
			}
		}
	}()

	return nil
}

func (this *AsyncProducer) Close() error {
	if this.producer == nil {
		return nil
	}

	close(this.chClose)
	return this.producer.Close()
}
