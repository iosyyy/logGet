package loggamet

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
)

type LogChanType struct {
	Value string
	Topic string
}
type LogConfType struct {
	Tail   *tail.Tail
	Topic  string
	Cancel context.CancelFunc
}

var (
	LogChannel = make(chan LogChanType, 50)

	ConfMap = make(map[string]*LogConfType)
)

func Init(filename string, topic string, ctx context.Context, cancel context.CancelFunc) (tails *tail.Tail) {
	config := tail.Config{
		ReOpen: true, // 重新打开
		Follow: true, // 是否跟随
		Location: &tail.SeekInfo{ // 从文件哪里开始读取
			Offset: 0,
			Whence: 2,
		},
		MustExist: false, // 文件不存在不报错
		Poll:      true,
	}
	var err error
	tails, err = tail.TailFile(filename, config)

	if err != nil {
		println(err)
		return
	}

	ConfMap[filename] = &LogConfType{
		Tail:   tails,
		Cancel: cancel,
		Topic:  topic,
	}
	go func() {

		var (
			line *tail.Line
			OK   bool
		)

		for true {
			select {
			case <-ctx.Done():
				tails.Cleanup()
				fmt.Println("监听通道: ", filename, "关闭")
				return
			case line, OK = <-tails.Lines:
				if !OK {
					fmt.Println("get file err: ", filename)
					return
				}
				LogChannel <- LogChanType{
					Value: line.Text,
					Topic: ConfMap[filename].Topic,
				}
			}
		}
	}()
	return
}
