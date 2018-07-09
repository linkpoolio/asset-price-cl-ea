package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"gopkg.in/guregu/null.v3"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
)

type Log struct {
	Message string
}

type Input struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type Output struct {
	Id        string   		    `json:"id"`
	Price     string   		    `json:"price"`
	Volume    string   		    `json:"volume"`
	Exchanges []string 		    `json:"exchanges"`
	Errors	  []*exchange.Error `json:"errors"`
}

type Params struct {
	Input
	Output
}

type RunResult struct {
	JobRunID     string        `json:"jobRunId"`
	Params       Params        `json:"data"`
	Status       string        `json:"status"`
	ErrorMessage null.String   `json:"error"`
	Pending      bool          `json:"pending"`
}

func Api() *rest.Api{
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/price", GetResponse),
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
