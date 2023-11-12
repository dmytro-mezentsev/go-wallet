package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Db DbConf `yaml:"db"`
}

type DbConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func GetConfig() Config {
	confContent, err := os.ReadFile("wallet/config/config.yaml")
	if err != nil {
		log.Fatal("can't read config file: ", err)
	}
	// expand environment variables
	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		log.Fatal("can't unmarshal config file: ", err)
	}
	return *conf
}
