package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
)

type Log struct {
	Message string
}

func Api() *rest.Api{
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/price/:base/:quote", GetResponse),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.SetPrefix("Asset Price CL Adaptor ")
	log.Print("Retrieving exchange trading pairs...")
	SetTradingPairs()
	log.Print("Set trading pairs, starting API...")

	api.SetApp(router)
	log.Print("Api started!")
	return api
}
