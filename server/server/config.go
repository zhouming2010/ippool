package server

import (
	"encoding/json"
	"os"
	"sync"
)

type DBServerConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"password"`
}

type IPServerConf struct {
	Port int `json:"port"`
}

type AppConf struct {
	DBServer DBServerConf `json:"db_server"`
	IPServer IPServerConf `json:"ip_server"`
}

var (
	instance *AppConf
	once     sync.Once
)

func GetAppConf() *AppConf {
	once.Do(func() {
		instance = &AppConf{}
	})
	return instance
}

func (h *AppConf) LoadConfig() {
	data, err := os.ReadFile("ippool_server.conf")
	if err != nil {
		panic("读取配置失败: " + err.Error())
	}
	err = json.Unmarshal(data, h)
	if err != nil {
		panic("解析配置失败: " + err.Error())
	}

}
