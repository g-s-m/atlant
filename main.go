package main

import (
	"atlant/requestor"
	_ "atlant/service"
	"fmt"
	"github.com/vrischmann/envconfig"
	"log"
)

type Config struct {
	Port string `envconfig:"PS_PORT,default=50000"`
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
