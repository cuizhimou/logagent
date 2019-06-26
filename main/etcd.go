package main

import (
	"awesomeProject/logagent/tailf"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strings"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
)

func initEtcd(addr string, key string) (collectConf [] tailf.CollectConf, err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr,},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("init error")
		return
	}
	logs.Debug("connect succ\n")
	//defer cli.Close()

	etcdClient = &EtcdClient{
		client: cli,
	}

	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}

	//var collectConf []tailf.CollectConf

	for _, ip := range localIpArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		etcdClient.keys = append(etcdClient.keys, etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			continue
		}
		cancel()
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				err = json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed, err:%v", err)
					continue
				}
				logs.Debug("log config is %v", collectConf)
			}
		}
	}

	initEtcdWatcher()

	return
}
func initEtcdWatcher() {
	for _, key := range etcdClient.keys {
		go watchKey(key)
	}
}

func watchKey(key string) {

	//初始化etcd客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.2.4:2379",},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("init error")
		return
	}

	//watch key
	for {
		rch := cli.Watch(context.Background(), key)
		var collectConf [] tailf.CollectConf
		var getConfSucc = true
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] config deleted", key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key[%s] unmarshal, err:%v", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc{
				tailf.UpdateConfig(collectConf)
			}
		}

	}
}
