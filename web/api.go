package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
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

	log.Print("Starting trading pairs ticker...")
	StartPairsTicker()
	log.Print("Set trading pairs, starting API...")

	api.SetApp(router)
	log.Print("API started!")
	return api
}
