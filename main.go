package main

import (
	"atlant/data/mongo"
	"atlant/service"
	"atlant/service/handler"
	"log"
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port     string `envconfig:"PS_ADDR_PORT,default=:50000"`
	Database string `envconfig:"DB_CONNECT,default=mongodb://127.0.0.1:27017"`
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

	repo := mongo.NewProductsActions(conf.Database, 60*time.Second)
	s := service.NewServer(handler.NewRequestHandler(repo))
	s.Run(conf.Port)
}
