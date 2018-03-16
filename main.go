package main

import (
	"github.com/linkpoolio/asset-price-cl-ea/web"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", web.Api().MakeHandler()))
}
