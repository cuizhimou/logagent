package main

import (
	"awesomeProject/logagent/tailf"
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
)

var (
	appConfig *Config
)

type Config struct {
	loglevel string
	logpath  string

	chanSize    int
	kafkaAddr   string
	collectConf [] tailf.CollectConf

	etcdAddr string
	etcdKey	 string
}

//根据配置收集日志
func loadCollecConf(conf config.Configer) (err error) {
	var cc tailf.CollectConf
	cc.LogPath = conf.String("collect::log_path")
	if len(cc.LogPath) == 0 {
		err = errors.New("invlid collect::log_path")
		return
	}

	cc.Topic = conf.String("collect::topic")
	if len(cc.LogPath) == 0 {
		err = errors.New("invlid collect::topic")
		return
	}

	appConfig.collectConf = append(appConfig.collectConf, cc)
	return
}

//加载配置文件中的配置
func loadConf(confType, filename string) (err error) {

	//初始化config对象
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	appConfig = &Config{}

	appConfig.loglevel = conf.String("logs::log_level")
	if len(appConfig.loglevel) == 0 {
		appConfig.loglevel = "debug"
	}

	appConfig.logpath = conf.String("logs::log_path")
	if len(appConfig.logpath) == 0 {
		appConfig.logpath = "./logs"
	}

	appConfig.chanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		appConfig.chanSize = 100
	}

	appConfig.kafkaAddr = conf.String("kafka::server_addr")
	if len(appConfig.kafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka address ")
		return
	}

	appConfig.etcdAddr = conf.String("etcd::addr")
	if len(appConfig.etcdAddr) == 0 {
		err = fmt.Errorf("invalid etcd address ")
		return
	}

	appConfig.etcdKey = conf.String("etcd::configKey")
	if len(appConfig.etcdKey) == 0 {
		err = fmt.Errorf("invalid etcd key ")
		return
	}

	err = loadCollecConf(conf)
	if err != nil {
		fmt.Printf("load collect conf failed ,err:%v\n", err)
		return
	}

	return
}
