package cmd

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/ghodss/yaml"
	"github.com/SevenIOT/windear/log"
	"github.com/SevenIOT/windear/server"
	"io/ioutil"
	"os"
	"time"
)

/**
 *
 * @author: schbook
 * @email: schbook@gmail.com
 * @date: 2018/4/12
 *
 */

const (
	Version               = "0.0.1"
	DefaultConfigFilePath = "config.yaml"
)

type Config struct {
	Redis struct {
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		Password        string `yaml:"password"`
		Database        int    `yaml:"database"`
		MaxIdle         int    `yaml:"maxIdle"`
		MaxTotal        int    `yaml:"maxTotal"`
		IdleTimeout     int    `yaml:"idleTimeout"`
		Timeout         int64  `yaml:"timeout"`
	} `yaml:"redis"`

	Host struct {
		IP string `yaml:"ip"`
	} `yaml:"host"`

	Rpc struct {
		Port int `yaml:"port"`
	} `yaml:"rpc"`

	Mqtt struct {
		Port int `yaml:"port"`
	} `yaml:"mqtt"`
}

var mConfig Config

func initConfig() bool {
	var configPath string
	help := flag.Bool("h", false, "this help")
	flag.StringVar(&configPath, "c", "", "config file path")
	// usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "windear version: %v\n\nUsage:\n", Version)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		return false
	}
	if configPath == "" {
		configPath = DefaultConfigFilePath
	}
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("open config file %v error:%v", configPath, err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &mConfig)
	if err != nil {
		log.Fatalf("config file parse error: %v", err.Error())
	}
	return true
}

func StartServer() {
	if initConfig() {
		redisPool := &redis.Pool{
			MaxIdle:         mConfig.Redis.MaxIdle,
			MaxActive:       mConfig.Redis.MaxTotal,
			IdleTimeout:     time.Duration(mConfig.Redis.IdleTimeout) * time.Millisecond,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", fmt.Sprintf("%v:%v", mConfig.Redis.Host, mConfig.Redis.Port),
					redis.DialConnectTimeout(time.Duration(mConfig.Redis.Timeout)*time.Millisecond),redis.DialPassword(mConfig.Redis.Password))
			},
		}
		server.NewServer(mConfig.Mqtt.Port, mConfig.Rpc.Port, redisPool, mConfig.Host.IP).Start()
	}
}
