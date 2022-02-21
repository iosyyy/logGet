package iosyyy

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var Producer sarama.SyncProducer

func Init(addr []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	/*msg := &sarama.ProducerMessage{
		Topic: "testtopic",
		Value: sarama.StringEncoder("hello world"),
	}*/
	var err error
	Producer, err = sarama.NewSyncProducer(addr, config)

	if err != nil {
		fmt.Println(err)
		return
	}

	/*	pid, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
		fmt.Printf("pid: %v, offset:%v", pid, offset)*/
}

func SendMessage(msg *sarama.ProducerMessage) (pid int32, offset int64) {
	pid, offset, err := Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return
}
