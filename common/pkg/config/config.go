package config

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/encoder/json"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/config/source/file"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Db struct {
	Address               string        `json:"address"`
	Database              string        `json:"database"`
	User                  string        `json:"user"`
	Password              string        `json:"password"`
	Charset               string        `json:"charset"`
	MaxConnections        int           `json:"max_connections"`
	MaxIdeConnections     int           `json:"max_ide_connections"`
	ConnectionMaxLifeTime time.Duration `json:"connection_max_life_time"`
}

type RabbitMq struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Configuration struct {
	Db Db `json:"db"`
	RabbitMq RabbitMq `json:"rabbit_mq"`
}

const etcdRootKey = "pg-service"

var serviceConfig *Configuration
var once sync.Once

func LoadConfig() *Configuration {
	once.Do(func() {
		//1.根据环境变量读取配置信息
		etcdConfigKey := os.Getenv("ETCD_CONFIG_KEY")
		serviceConfig = &Configuration{}
		encoder := json.NewEncoder()
		fmt.Println(os.Getwd())
		fileSource := file.NewSource(file.WithPath("./config.json"), source.WithEncoder(encoder))
		etcdSource := etcd.NewSource(
			etcd.WithAddress(strings.Split(os.Getenv("MICRO_REGISTRY_ADDRESS"), ",")[0]),
			etcd.WithPrefix(etcdRootKey+"/"+etcdConfigKey),
			etcd.StripPrefix(false),
			source.WithEncoder(encoder),
		)

		//2.根据环境变量读取etcd或本地json
		conf, _ := config.NewConfig()
		var err error
		if os.Getenv("ENABLE_REMOTE_CONFIG") == "true" {
			err = conf.Load(
				fileSource, // 将文件配置作为默认值
				etcdSource, // 会覆盖上面的文件配置
			)
		} else {
			err = conf.Load(fileSource)
		}
		if err != nil {
			// 加载数据源失败
			log.Fatalf("load config fail: %v", err)
		}

		//3.加载配置文件
		err = conf.Get(etcdRootKey, etcdConfigKey).Scan(serviceConfig)
		if err != nil {
			// 读取远程配置失败
			log.Fatalf("etcd load config fail: %v", err)
		}
		log.Printf("config：%v", serviceConfig)
		log.Printf("config map：%v", conf.Map())

		//4.开启协程监听配置变更
		go func() {
			for {
				time.Sleep(time.Second * 5) // delay after each request

				w, err := conf.Watch(etcdRootKey, etcdConfigKey)
				if err != nil {
					log.Printf("wactch config fail: %v", err)
					continue
				}

				// wait for next value
				value, err := w.Next()
				if err != nil {
					log.Printf("wacth cofig fail: %v", err)
					continue
				}

				value.Scan(serviceConfig)
				log.Printf("item change：%s", serviceConfig)
			}
		}()
	})

	return serviceConfig
}
