package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
	"fmt"
	"sync"
	"net/http"
)

func GetPrice(w rest.ResponseWriter, r *rest.Request) {
	base := r.PathParam("base")
	quote := r.PathParam("quote")

	// Call the exchanges through a concurrently
	responses, errors := getExchangeResponses(base, quote)
	if len(errors) > 0 {
		rest.Error(w, fmt.Sprintf("Following errors were given: %s", errors), http.StatusInternalServerError)
		return
	}

	// Calculate the weighted average based on volume
	var totalVolume float64
	for _, response := range responses {
		totalVolume += response.Volume
	}
	var weightedPrice float64
	for _, response := range responses {
		weightedPrice += (response.Volume / totalVolume) * response.Price
	}
	w.WriteJson(&exchange.Response{fmt.Sprintf("%s-%s", base, quote), weightedPrice, totalVolume})
}

func getExchangeResponses(base, quote string) ([]*exchange.Response, []error) {
	supportedExchanges := exchange.GetSupportedExchanges()

	var wg sync.WaitGroup
	wg.Add(len(supportedExchanges))

	var responses []*exchange.Response
	var errors []error

	for _, exc := range supportedExchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			response, e := exc.GetResponse(base, quote)
			if e != nil {
				error := fmt.Errorf("error retrieving information from %s", exc.GetConfig().Name)
				errors = append(errors, error)
			}
			responses = append(responses, response)
		}(exc)
	}

	wg.Wait()

	return responses, errors
}