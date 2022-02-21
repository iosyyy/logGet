package conf

type AppConf struct {
	KafkaConf `ini:"kafka"`
	LogConf   `ini:"log"`
	EtcdConf  `ini:"etcd"`
}
type KafkaConf struct {
	Topic   string `ini:"topic"`
	Address string `ini:"address"`
	Key     string `ini:"key"`
}

type EtcdConf struct {
	Timeout int    `ini:"timeout"`
	Address string `ini:"address"`
}

type LogConf struct {
	Filename string `ini:"filename"`
	Key      string `ini:"key"`
}
