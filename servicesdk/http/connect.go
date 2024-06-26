package http

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"git.yuetanggame.com/zfish/fishpkg/logs"
	. "git.yuetanggame.com/zfish/fishpkg/servicesdk/core/helper"
	p "git.yuetanggame.com/zfish/fishpkg/sprotocol/http"
)

type _GWStatus uint8

const (
	_GW_STATUS_IDLE _GWStatus = iota
	_GW_STATUS_PENDING
	_GW_STATUS_RUNNING
)

type _GWList struct {
	sync.RWMutex
	m map[string]*_GWContext
}

type _GWContext struct {
	so       *p.Socket
	status   _GWStatus
	cancelFn context.CancelFunc
}

var gwList _GWList = _GWList{m: make(map[string]*_GWContext, 0)}

func (this *_GWList) GetReadyCount() int {
	this.RLock()
	defer this.RUnlock()

	return len(this.m)
}

func (this *_GWList) GetAll() []*p.Socket {
	this.RLock()
	defer this.RUnlock()

	var conns []*p.Socket
	for _, v := range this.m {
		if v.status == _GW_STATUS_RUNNING {
			conns = append(conns, v.so)
		}
	}

	return conns
}

func (this *_GWList) Del(key string) {
	this.Lock()
	this.Unlock()

	gw, ok := this.m[key]
	if !ok {
		return
	}

	if gw.status == _GW_STATUS_PENDING {
		gw.cancelFn()
	}

	if gw.status == _GW_STATUS_RUNNING {
		gw.so.Close()
	}

	delete(this.m, key)
}

func (this *_GWList) Add(key string) {
	this.Lock()
	defer this.Unlock()

	gw, ok := this.m[key]
	if !ok {
		gw = &_GWContext{}
		this.m[key] = gw
	}

	if gw.status == _GW_STATUS_PENDING {
		return
	}

	if gw.status == _GW_STATUS_RUNNING {
		gw.so.SetContext(nil)
		gw.so.Close()
	}

	ctx, fn := context.WithCancel(context.TODO())
	gw.cancelFn = fn
	gw.status = _GW_STATUS_PENDING

	go this.dial(ctx, key)
}

func (this *_GWList) dial(ctx context.Context, key string) {
	for {
		select {
		case <-ctx.Done():
			logs.Infof("Connect to gateway[%v] canceled", key)
			return
		default:
			cli, err := net.Dial("tcp", key)
			if err != nil {
				logs.Errorf("Connect to gateway[%v] err: %v", key, err.Error())
				time.Sleep(5 * time.Second)
				continue
			}

			logs.Debugf("Connect to gateway[%v] success", key)

			so := p.NewSocket(cli,
				p.SetFilter(&Filter{}),
				p.SetNotify(&Notify{}),
				p.SetTimeout(srv.Timeout()))

			so.Start()

			if err := srv.register(so); err != nil {
				logs.Errorf("Register err: %s", err.Error())

				if strings.Contains(err.Error(), "send timeout") {
					this.Lock()
					if gw, ok := this.m[key]; ok {
						gw.status = _GW_STATUS_IDLE
					}
					this.Unlock()
				}

				so.Close()
				return
			}

			this.Lock()
			defer this.Unlock()

			gw, ok := this.m[key]
			if !ok {
				so.Close()
				return
			}

			if gw.status != _GW_STATUS_PENDING {
				so.Close()
				return
			}

			gw.so = so
			gw.cancelFn = nil
			gw.status = _GW_STATUS_RUNNING

			so.SetContext(key)

			return
		}
	}
}

func dialGateway(ips []string) error {
	if len(ips) == 0 {
		return errors.New("gateway addr isn't specified")
	}
	for _, k := range ips {
		if HostAddrCheck(k) == false {
			logs.Errorf("Bad gateway address: %s", k)
			continue
		}

		gwList.Add(k)
		logs.Debugf("Add gateway[%v] to pool", k)
	}

	return nil
}

func fetchGateway() error {
	var (
		keys []string
		err  error
	)

	pfx := srv.GatewayDir()

	if keys, _, err = srv.EtcdConn().GetKvWithPrefix(pfx); err != nil {
		return err
	}

	for _, k := range keys {
		key := strings.TrimPrefix(k, pfx)
		if HostAddrCheck(key) == false {
			logs.Errorf("Bad gateway address: %s", key)
			continue
		}

		gwList.Add(key)
		logs.Debugf("Add gateway[%v] to pool", key)
	}

	monitor := func() {
		for {
			time.Sleep(3 * time.Second)
			if gwList.GetReadyCount() == 0 {
				logs.Waringf("No gateway is available ...")
			}
		}
	}

	go monitor()

	go watchGateway()

	return nil
}

func watchGateway() {
	pfx := srv.GatewayDir()

	for {
		time.Sleep(1 * time.Second)
		rch := srv.EtcdConn().WatchKeyWithPrefix(pfx)
		logs.Debugf("Wath: %v", pfx)
		for wresp := range rch {
			if wresp.Canceled {
				logs.Infof("Watch %v canceled", pfx)
				break
			}

			if err := wresp.Err(); err != nil {
				logs.Errorf("Watch %v err: %s", pfx, err.Error())
				break
			}

			for _, ev := range wresp.Events {
				switch ev.Type.String() {
				case "PUT":
					logs.Infof("Add ETCD [key:%s,value:%s]", string(ev.Kv.Key), string(ev.Kv.Value))
					key := strings.TrimPrefix(string(ev.Kv.Key), pfx)

					if HostAddrCheck(key) == false {
						logs.Errorf("Bad gateway address: %s", key)
						continue
					}

					gwList.Add(key)
					logs.Debugf("Add gateway[%v] to pool", key)
				case "DELETE":
					logs.Infof("Delete ETCD [key:%s,value:%s]", string(ev.Kv.Key), string(ev.Kv.Value))
					key := strings.TrimPrefix(string(ev.Kv.Key), pfx)

					if HostAddrCheck(key) == false {
						logs.Errorf("Bad gateway address: %s", key)
						continue
					}

					gwList.Del(key)
					logs.Debugf("Delete gateway[%v] from pool", key)
				}
			}
		}
	}
}
