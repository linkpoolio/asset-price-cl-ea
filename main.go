package main

import (
	"fmt"
	"github.com/linkpoolio/asset-price-cl-ea/web"
	"log"
	"net/http"
)

func main() {
	web.InitialiseConfig()

	log.Print("Chainlink Asset Price Adaptor")
	log.Printf("Starting to serve on port %d", web.Config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", web.Config.Port), web.Api().MakeHandler()))
}
