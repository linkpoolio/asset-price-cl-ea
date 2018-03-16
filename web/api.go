package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
)

func Api() *rest.Api{
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/price/:base/:quote", GetPrice),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api
}
