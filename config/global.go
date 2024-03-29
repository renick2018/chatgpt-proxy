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
	GPTServers []struct {
		Plus     bool   `yaml:"plus" json:"plus"`
		Host     string `yaml:"host"`
		Email    string `yaml:"email" json:"email"`
		Password string `yaml:"password" json:"password"`
		ApiKey   string `yaml:"apiKey" json:"api_key"`
	} `yaml:"GPTServers"`                // v2
	Emails      []string `yaml:"Emails"` // alert emails
	EmailServer struct {
		Sender   string `yaml:"sender"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"EmailServer"`
}

func init() {
	path, _ := os.Getwd()
	fileData, err := os.ReadFile(fmt.Sprintf("%s/.conf.yml", path))
	if err != nil {
		logger.Error(fmt.Sprintf("load conf file error: %+v", err))
		return
	}

	if e := yaml.Unmarshal(fileData, &Global); e != nil {
		logger.Error(fmt.Sprintf("unmarshal conf file error: %+v", err))
	}
}
