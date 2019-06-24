package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
}

var (
	etcdClient *EtcdClient
)

func initEtcd(addr string, key string) (err error) {
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
	for _, ip := range localIpArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			continue
		}
		cancel()
		for k, v := range resp.Kvs {
			fmt.Println(k,v)
		}
	}
	return
}
