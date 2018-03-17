package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
	"fmt"
	"sync"
	"net/http"
)

type Error struct {
	Error string
	StatusCode int
	Errors []*exchange.Error
}

func GetPrice(w rest.ResponseWriter, r *rest.Request) {
	base := r.PathParam("base")
	quote := r.PathParam("quote")

	// Call the exchanges concurrently
	responses, errors := getExchangeResponses(base, quote)
	if len(errors) > 0 {
		errorObj := &Error{
			fmt.Sprintf("errors given when getting prices from exchanges"),
			http.StatusInternalServerError,
			errors,
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteJson(&errorObj)
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

func getExchangeResponses(base, quote string) ([]*exchange.Response, []*exchange.Error) {
	supportedExchanges := exchange.GetSupportedExchanges()

	var wg sync.WaitGroup
	wg.Add(len(supportedExchanges))

	var responses []*exchange.Response
	var errors []*exchange.Error

	for _, exc := range supportedExchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			response, err := exc.GetResponse(base, quote)
			if err != nil {
				errors = append(errors, err)
			}
			responses = append(responses, response)
		}(exc)
	}

	wg.Wait()

	return responses, errors
}