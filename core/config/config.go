package config

import (
	"github.com/spf13/viper"
	"sync"
)

type OnceFlag struct {
	IsInit bool //判断配置是否初始化，已初始化为true
	sync.Mutex
}

var configs = map[string]*viper.Viper{}

func LoadConfigFile(name string, path string, configType string) {
	c := viper.New()
	c.SetConfigName(name)
	c.SetConfigType(configType)
	c.AddConfigPath(path)

	if err := c.ReadInConfig(); err != nil {
		panic("Config File(" + path + "/" + name + "." + configType + ") parse error:" + err.Error())
	}

	configs[name] = c
}

func ReadConfigOnce(config string, key string, r interface{}, f *OnceFlag) {
	if f.IsInit {
		return
	}
	f.Lock()
	defer f.Unlock()
	if f.IsInit {
		return
	}

	ReadConfig(config, key, r)
	f.IsInit = true
}

func ReadConfig(config string, key string, r interface{}) {
	c, ok := configs[config]
	if !ok {
		panic("cannot find config " + config)
	}
	c.UnmarshalKey(key, r)
}
