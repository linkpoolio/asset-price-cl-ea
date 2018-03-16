package main

import (
	"log"
	"net/http"
	"github.com/linkpoolio/asset-price-cl-ea/web"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", web.Api().MakeHandler()))
}