package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/kkkkiven/fishpkg/logs"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	defReqTimeout  = 2 * time.Second
	defDialTimeout = 5 * time.Second
)

type Client struct {
	obj *clientv3.Client
}

// New 创建连接
func New(addrs []string, user, pass string, timeout ...time.Duration) (*Client, error) {
	dialTimeout := defDialTimeout
	if len(timeout) > 0 {
		dialTimeout = timeout[0]
	}

	object, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: dialTimeout,
		Username:    user,
		Password:    pass,
	})

	if err != nil {
		return nil, err
	}

	return &Client{object}, nil
}

func (cli *Client) GetClient() *clientv3.Client {
	return cli.obj
}

func (cli *Client) Close() error {
	return cli.obj.Close()
}

func (cli *Client) Get(key string, timeout ...time.Duration) (string, error) {
	reqTimeout := defReqTimeout
	if len(timeout) > 0 {
		reqTimeout = timeout[0]
	}

	ctx, _ := context.WithTimeout(context.TODO(), reqTimeout)

	rsp, err := cli.obj.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(rsp.Kvs) > 0 {
		return string(rsp.Kvs[0].Value), nil
	}

	return "", nil
}

func (cli *Client) Set(key, value string, timeout ...time.Duration) error {
	reqTimeout := defReqTimeout
	if len(timeout) > 0 {
		reqTimeout = timeout[0]
	}

	ctx, _ := context.WithTimeout(context.TODO(), reqTimeout)
	_, err := cli.obj.Put(ctx, key, value)

	return err
}

func (cli *Client) Del(key string, timeout ...time.Duration) error {
	reqTimeout := defReqTimeout
	if len(timeout) > 0 {
		reqTimeout = timeout[0]
	}

	ctx, _ := context.WithTimeout(context.TODO(), reqTimeout)
	_, err := cli.obj.Delete(ctx, key)

	return err
}

func (cli *Client) DelKeysWithPrefix(prefix string, timeout ...time.Duration) error {
	reqTimeout := defReqTimeout
	if len(timeout) > 0 {
		reqTimeout = timeout[0]
	}

	ctx, _ := context.WithTimeout(context.TODO(), reqTimeout)
	_, err := cli.obj.Delete(ctx, prefix, clientv3.WithPrefix())

	return err
}

func (cli *Client) GetKvWithPrefix(prefix string, timeout ...time.Duration) ([]string, []string, error) {
	reqTimeout := defReqTimeout
	if len(timeout) > 0 {
		reqTimeout = timeout[0]
	}

	ctx, _ := context.WithTimeout(context.TODO(), reqTimeout)
	rsp, err := cli.obj.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, nil, err
	}

	var keys, values []string
	for _, ev := range rsp.Kvs {
		keys = append(keys, string(ev.Key))
		values = append(values, string(ev.Value))
	}

	return keys, values, nil
}

func (cli *Client) WatchKey(key string) clientv3.WatchChan {
	return cli.obj.Watch(context.TODO(), key)
}

func (cli *Client) WatchKeyWithPrefix(prefix string) clientv3.WatchChan {
	return cli.obj.Watch(context.TODO(), prefix, clientv3.WithPrefix())
}

func (cli *Client) SetKvAndKeepAlive(key, value string, ttl, interval int64) error {

	go func() {
		logs.Debugf("Keepalive key: %v, value: %v", key, value)

		for {
		LOOP:
			reqTimeout := defReqTimeout

			ctx, _ := context.WithTimeout(context.TODO(), reqTimeout)
			rsp, err := cli.obj.Grant(ctx, ttl)
			if err != nil {
				logs.Errorf("Keepalive err: %v", err.Error())
				continue
			}

			ctx, _ = context.WithTimeout(context.TODO(), reqTimeout)
			if _, err := cli.obj.Put(ctx, key, value, clientv3.WithLease(rsp.ID)); err != nil {
				logs.Errorf("Keepalive err: %v", err.Error())
				continue
			}

			ticker := time.NewTicker(time.Duration(interval) * time.Second)
			for {
				select {
				case <-ticker.C:
					if _, err := cli.obj.KeepAliveOnce(context.TODO(), rsp.ID); err != nil {
						logs.Errorf("Keepalive key:%s, err: %v", key, err.Error())
						if ok := strings.Contains(err.Error(), "requested lease not found"); ok {
							logs.Debugf("Rekeepalive key: %v, value: %v", key, value)
							ticker.Stop()
							goto LOOP
						}
						continue
					}
				}
			}
		}
	}()

	return nil
}
