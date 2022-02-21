package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

var Client *clientV3.Client
var typeOfConf = &TypeConf{}
var logListenConf = &LogListenConf{}

type TypeConf struct {
	Key string

	Value []string
}

type LogListen struct {
	Filename string `json:"filename"`
	Topic    string `json:"topic"`
}

type LogListenConf struct {
	Key string

	Value []LogListen
}

func Init(address string, timeout time.Duration) {
	var err error
	Client, err = clientV3.New(clientV3.Config{
		Endpoints:   []string{address},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Watch(key string, callback func([]byte)) {
	go func() {
		watch := Client.Watch(context.TODO(), key)

		for response := range watch {
			for _, event := range response.Events {
				callback(event.Kv.Value)
			}
		}
	}()
}

func GetConfByKey(ctx context.Context, key string) *TypeConf {
	value, err := Client.Get(ctx, key)
	if err != nil {
		return nil
	}
	for _, val := range value.Kvs {
		typeOfConf.Key = string(val.Key)
		err := json.Unmarshal(val.Value, &typeOfConf.Value)
		fmt.Println(typeOfConf.Value)
		if err != nil {
			return nil
		}
	}
	return typeOfConf
}

func SetListenConfValue(value []byte) *LogListenConf {
	err := json.Unmarshal(value, &logListenConf.Value)
	if err != nil {
		fmt.Println("err: ", err)
		return nil
	}

	return logListenConf

}

func GetListenLogConfByKey(ctx context.Context, key string) *LogListenConf {
	value, err := Client.Get(ctx, key)
	if err != nil {
		return nil
	}
	for _, val := range value.Kvs {
		logListenConf.Key = string(val.Key)
		err := json.Unmarshal(val.Value, &logListenConf.Value)
		if err != nil {
			return nil
		}
	}
	return logListenConf
}

/*func getEtcd() (err error) {
/*client, err := clientV3.New(clientV3.Config{
	Endpoints:   []string{"127.0.0.1:2379"},
	DialTimeout: 5 * time.Second,
})*/
/*go func() {
		watch := client.Watch(context.TODO(), "hello")

		for response := range watch {
			for _, value := range response.Events {
				fmt.Printf("\nXvalue: %s", value.Kv.Value)
			}
		}
	}()

	if err != nil {
		return
	}
	defer func(client *clientV3.Client) {
		err := client.Close()
		if err != nil {
			fmt.Println("err: ", err)
		}
	}(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, "hello", "world")
	cancel()
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, "hello", clientV3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}
	for _, val := range resp.Kvs {
		fmt.Printf("key: %s value: %s", val.Key, val.Value)

	}
	return nil
}*/
