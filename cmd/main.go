// Package main TutuPlapak Product Purchase API
//
//	@title			TutuPlapak Product Purchase API
//	@version		1.0
//	@description	A product purchase API for TutuPlapak
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8001
//	@BasePath	/api/v1
//
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
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
