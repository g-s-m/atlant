package main

import (
	"atlant/service"
	"log"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port string `envconfig:"PS_ADDR_PORT,default=:50000"`
}

func getConfig() *Config {
	var conf Config
	if err := envconfig.Init(&conf); err != nil {
		log.Fatalf("Can't read config: %v", err)
	}
	return &conf
}

func main() {
	conf := getConfig()
	service.RunService(conf.Port)
}
