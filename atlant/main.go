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
	Port   string `envconfig:"PS_ADDR_PORT,default=:50000"`
	DbAddr string `envconfig:"DB_CONNECT,default=cluster0.9aypu.mongodb.net"`
	DbPswd string `envconfig:"DB_PSWD"`
	DbUser string `envconfig:"DB_USER,default=auser"`
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

	repo := mongo.NewProductsActions(conf.DbAddr, conf.DbUser, conf.DbPswd, 60*time.Second)
	s := service.NewServer(handler.NewRequestHandler(repo))
	s.Run(conf.Port)
}
