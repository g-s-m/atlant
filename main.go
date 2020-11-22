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

	repo := mongo.NewProductsActions(conf.Database, 10*time.Second)
	repo.Save("product_1", 1.12)
	repo.Save("product_4", 2.12)
	repo.Save("product_2", 6.12)
	repo.Save("product_3", 4.12)

	r, err := repo.LoadByChangeCount(0, 2, false)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	for _, e := range r {
		log.Printf(" %v ", *e)
	}

	s := service.NewServer(handler.NewRequestHandler(repo))
	s.Run(conf.Port)

	//repo := data.CreateRepo(db)
	//repo.LoadByProduct()
}
