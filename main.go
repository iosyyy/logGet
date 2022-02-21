package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	clientV3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/ini.v1"
	"logGet/src/github.com/conf"
	"logGet/src/github.com/etcd"
	"logGet/src/github.com/iosyyy"
	"logGet/src/github.com/loggamet"
	"time"
)

var appConf = &conf.AppConf{}

func run() {

	for messageStr := range loggamet.LogChannel {
		fmt.Println("发送消息: ", messageStr.Value)
		fmt.Println("接收的topic: ", messageStr.Topic)
		pid, offset := iosyyy.SendMessage(&sarama.ProducerMessage{
			Topic: messageStr.Topic,
			Value: sarama.StringEncoder(messageStr.Value),
		})

		fmt.Printf("发送消息成功 pid: %v, offset: %v\n", pid, offset)

	}

}

func changeLog(value []byte) {
	confValue := etcd.SetListenConfValue(value)
	listens := confValue.Value
	m := map[string]string{}
	for _, listen := range listens {
		m[listen.Filename] = listen.Topic
	}
	for filename := range loggamet.ConfMap {
		topic, ok := m[filename]

		if !ok {
			loggamet.ConfMap[filename].Cancel()
			delete(loggamet.ConfMap, filename)

		} else {
			fmt.Printf("监听%v的topic从%v -> %v\n", filename, loggamet.ConfMap[filename].Topic, topic)

			loggamet.ConfMap[filename].Topic = topic
		}
	}

	for filename, topic := range m {
		_, ok := loggamet.ConfMap[filename]
		if !ok {
			ctx, cancel := context.WithCancel(context.Background())
			loggamet.Init(filename, topic, ctx, cancel)
		}
	}
}

func main() {
	err := ini.MapTo(appConf, "./src/resourses/config.ini")
	if err != nil {
		return
	}
	etcd.Init(appConf.EtcdConf.Address, time.Duration(appConf.EtcdConf.Timeout)*time.Second)
	etcd.Watch(appConf.LogConf.Key, changeLog)
	defer func(EtcdClient *clientV3.Client) {
		err := EtcdClient.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(etcd.Client)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	kafkaConf := etcd.GetConfByKey(ctx, appConf.KafkaConf.Key)
	cancel()
	iosyyy.Init(kafkaConf.Value)
	defer func(Producer sarama.SyncProducer) {
		err := Producer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(iosyyy.Producer)
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	logConf := etcd.GetListenLogConfByKey(ctx, appConf.LogConf.Key)
	cancel()
	for _, value := range logConf.Value {
		ctx, cancel := context.WithCancel(context.Background())
		loggamet.Init(value.Filename, value.Topic, ctx, cancel)
	}
	run()
}
