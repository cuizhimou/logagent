package main

import (
	"awesomeProject/logagent/kafka"
	"awesomeProject/logagent/tailf"
	"github.com/astaxie/beego/logs"
	"time"
)

func serverRun() (err error) {
	for {
		msg := tailf.GetOneLine()
		err := sendToKafka(msg)
		if err != nil {
			logs.Error("send to kafaka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

func sendToKafka(msg *tailf.TextMsg)(err error)  {
	//logs.Debug("read mgs:%s, topic:%s",msg.Msg,msg.Topic)
	err = kafka.SendToKafka(msg.Msg,msg.Topic)
	return
}
