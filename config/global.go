package config

import (
	"chatgpt-proxy/lib/logger"
	"fmt"
	"os"
)
import "gopkg.in/yaml.v2"

var Global struct {
	ApiSalt         string   `yaml:"ApiSalt"`
	ChatServerAddrs []string `yaml:"ChatServerAddrs"` // v1
	GPTServers      []struct {
		Host     string `yaml:"host"`
		Email    string `yaml:"email" json:"email"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"GPTServers"` // v2
	Emails      []string `yaml:"Emails"` // alert emails
	EmailServer struct {
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
		logger.Error(fmt.Sprintf("unmarshal conf file error: %+v", err))
	}
}
