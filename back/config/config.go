package config

import (
	"application-web/global"
	"application-web/logger"
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	configurePath = "/etc/application-web/config" // 配置文件所在目录
)

var Conf Config

func LoadConfig() {
	Conf = Config{}
	Conf.v = viper.New()

	v := Conf.Application.v
	v.SetConfigFile("/etc/application-web/config/application-web.yaml")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil { // viper解析配置文件
		fmt.Printf("load config path: %s, error: %s", configurePath, err)
	}

	if err := v.Unmarshal(&global.System); err != nil {
		fmt.Printf("Unable to decode into struct, %s", err)
	}

}

type Application struct {
	name string
	v    *viper.Viper
}

func (c *Application) GetName() string {
	return c.name
}

func (c *Application) GetViper() *viper.Viper {
	return c.v
}

type Config struct {
	Application
	rwMutex sync.RWMutex
}

func (sc *Config) RLock() {
	sc.rwMutex.RLock()
}

func (sc *Config) RUnlock() {
	sc.rwMutex.RUnlock()
}

func (sc *Config) Lock() {
	sc.rwMutex.Lock()
}

func (sc *Config) Unlock() {
	sc.rwMutex.Unlock()
}

// 监控配置文件变化并热加载程序
func watchConfig(v *viper.Viper, callback func(in fsnotify.Event)) {
	v.OnConfigChange(callback)
	v.WatchConfig()
}

func SaveConfig() error {
	out, err := yaml.Marshal(global.System)
	if err != nil {
		logger.Error("error marshalling yaml: %v", err)
		return err
	}

	err = os.WriteFile("/etc/application-web/config/application-web.yaml", out, 0644) // 0644 是文件权限
	if err != nil {
		logger.Error("error marshalling yaml: %v", err)
		return err
	}

	return nil
}
