package main

import (
	"awesomeProject/logagent/kafka"
	"awesomeProject/logagent/tailf"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func main() {

	//读取加载配置文件
	filename := "./conf/logagent.conf"
	err := loadConf("ini", filename)
	if err != nil {
		fmt.Printf("load conf faild, err:%v\n", err)
		panic("load conf failed")
		return
	}

	//初始化日志
	err = initLogger()
	if err != nil {
		fmt.Printf("load logger failed, err:%v\n", err)
		panic("load logger failed")
		return
	}
	logs.Debug("load conf succ, config:%v", appConfig)

	//初始化Etcd、获取配置
	collectConf, err := initEtcd(appConfig.etcdAddr, appConfig.etcdKey)
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}
	logs.Debug("initialize etcd succ")

	//初始化tailconfig
	err = tailf.InitTail(collectConf, appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed, err:%v", err)
		return
	}
	logs.Debug("initialize tailf succ")

	//初始化kafka配置
	err = kafka.InitKafka(appConfig.kafkaAddr)
	if err != nil {
		logs.Error("initkafka failed, err:%v", err)
	}

	logs.Debug("initialize all succ")

	err = serverRun()
	if err != nil {
		logs.Error("serverRun failed, err:%v", err)
		return
	}
	logs.Info("program exited")

}
