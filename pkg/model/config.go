package model

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	configPath = "./config/config.yaml"
)

var GlobalConfig *Config

type Config struct {
	Accounts []Account `yaml:"accounts" json:"accounts"`
}

type Account struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

func GetConfig() *Config {
	if GlobalConfig != nil {
		return GlobalConfig
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Error reading config file:", err)
		return nil
	}
	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		log.Println("Error parsing config file:", err)
		return nil
	}
	return GlobalConfig
}
