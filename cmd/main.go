package main

import (
	"log"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/cmd/app"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log.Println("start running app")

	app.Run(cfg)
}
