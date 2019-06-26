package main

import (
	"awesomeProject/logagent/tailf"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
	"context"
)


const (
	EtcdKey = "/oldboy/backend/logagent/config/192.168.0.30"
)




func SetLogConfToEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.2.4:2379",},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("init error")
	}
	fmt.Printf("connect succ\n")
	defer cli.Close()

	var logConfArr []tailf.CollectConf
	logConfArr = append(logConfArr, tailf.CollectConf{
		LogPath:  "/Users/cui/tmp/access.log",
		Topic: "nginx_log",
	})
	logConfArr = append(logConfArr, tailf.CollectConf{
		LogPath:  "/Users/cui/tmp/error.log",
		Topic: "nginx_log_err",
	})

	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey,string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}


	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx,EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s,:%s\n", ev.Key, ev.Value)
	}

}

func main()  {
	SetLogConfToEtcd()
}

