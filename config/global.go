package config

import (
	"chatgpt-proxy/lib/logger"
	"fmt"
	"os"
)
import "gopkg.in/yaml.v2"

var Global struct {
	ApiSalt         string   `yaml:"ApiSalt"`
	ChatServerAddrs []string `yaml:"ChatServerAddrs"`
	Emails          []string `yaml:"Emails"`
	EmailServer     struct {
		Sender   string `yaml:"sender"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"EmailServer"`
}

func init() {
	fileData, err := os.ReadFile("./.conf.yml")
	if err != nil {
		logger.Error(fmt.Sprintf("load conf file error: %+v", err))
		return
	}

	if e := yaml.Unmarshal(fileData, &Global); e != nil {
		logger.Error(fmt.Sprintf("load conf file error: %+v", err))
	}
}
