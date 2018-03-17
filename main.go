package main

import (
	"github.com/linkpoolio/asset-price-cl-ea/web"
	"log"
	"net/http"
	"os"
	"fmt"
)

func main() {
	port := ""
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), web.Api().MakeHandler()))
}
