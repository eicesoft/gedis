package config

import (
	"github.com/spf13/viper"
	"sync"
)

var (
	config *Config   //服务配置
	once   sync.Once //单例实现
)

type Config struct {
	Server struct {
		Password  string `toml:"Password"`  //服务器密码
		Port      string `toml:"Port"`      //服务器默认监听端口
		EnableAof bool   `toml:"EnableAof"` //是否开启Aof
	} `toml:"Server"`
}

func Get() *Config {
	once.Do(func() {
		config = &Config{} //viper config 实现
		viper.SetConfigName("gedis")
		viper.SetConfigType("toml")
		viper.AddConfigPath("config")

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})

	return config
}
