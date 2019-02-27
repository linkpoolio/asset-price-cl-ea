package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

func Api() *rest.Api{
	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)
	router, err := rest.MakeRouter(
		rest.Post("/price", GetResponse),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("starting trading pairs ticker")
	StartPairsTicker()
	log.Print("set trading pairs, starting api")

	api.SetApp(router)
	log.Print("api started")
	return api
}
