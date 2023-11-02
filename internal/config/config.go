package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Storage struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Database   string `yaml:"database"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		DriverName string `yaml:"driverName"`
	} `yaml:"storage"`
	Nats struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Cluster_id string `yaml:"cluster_id"`
		Client_id  string `yaml:"client_id"`
		Channel    string `yaml:"channel"`
	} `yaml:"nats-server"`
	Http struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http-server"`
}

func (cfg *Config) InitFile() {
	f, err := os.Open("internal/config/config.yml")
	if err != nil {
		log.Println("Can't open config.yml")
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
