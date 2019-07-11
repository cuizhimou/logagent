package tailf

import (
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

//一个文件对象

type TailObj struct {
	tail     *tail.Tail
	conf     CollectConf
	status   int
	exitChan chan int
}

type TextMsg struct {
	Msg   string
	Topic string
}

//文件对象数组
type TailObjMgr struct {
	tailobjs [] *TailObj
	msgChan  chan *TextMsg
	lock sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

//更新配置

func UpdateConfig(confs []CollectConf) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
	for _, oneConf := range confs {
		var isRunning = false
		for _, obj := range tailObjMgr.tailobjs {
			if oneConf.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}
		if isRunning {
			continue
		}
		createNewTask(oneConf)
	}

	var tailobjs []*TailObj
	for _, obj := range tailObjMgr.tailobjs {
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailobjs = append(tailobjs,obj)
	}
	tailObjMgr.tailobjs=tailobjs
	return
}

func createNewTask(conf CollectConf) {
	obj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}
	tails, tailerr := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
	})
	if tailerr != nil {
		logs.Error("collect filename[%s] failed, err:%v", conf.LogPath, tailerr)
		return
	}

	obj.tail = tails
	tailObjMgr.tailobjs = append(tailObjMgr.tailobjs, obj)

	go readFromTail(obj)

}

//初始化tail
func InitTail(conf []CollectConf, chanSize int) (err error) {
	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}

	if len(conf) == 0 {
		logs.Error("invalid config for log_collect, conf:%v", conf)
		return
	}

	for _, v := range conf {
		createNewTask(v)
	}

	return
}

//读日志
func readFromTail(tailObj *TailObj) {
	for {
		select {
		case line, ok := <-tailObj.tail.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}
			tailObjMgr.msgChan <- textMsg
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited, conf:%v", tailObj.conf)
			return
		}
	}
}