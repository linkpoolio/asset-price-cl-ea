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

type Response struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
	Volume float64 `json:"volume"`
	Exchanges []string `json:"exchanges"`
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
	} else if len(responses) == 0 {
		errorObj := &Error{
			fmt.Sprintf("no exchanges support that trading pair"),
			http.StatusBadRequest,
			nil,
		}
		w.WriteHeader(http.StatusBadRequest)
		w.WriteJson(&errorObj)
		return
	}

	// Calculate the weighted average based on volume
	resp := &Response{Id: fmt.Sprintf("%s-%s", base, quote)}
	for _, response := range responses {
		resp.Exchanges = append(resp.Exchanges, response.Name)
		resp.Volume += response.Volume
	}
	for _, response := range responses {
		resp.Price += (response.Volume / resp.Volume) * response.Price
	}
	w.WriteJson(resp)
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
			if isSupportedExchange(exc, base, quote) {
				response, err := exc.GetResponse(base, quote)
				if err != nil {
					errors = append(errors, err)
				}
				responses = append(responses, response)
			}
		}(exc)
	}

	wg.Wait()

	return responses, errors
}

func isSupportedExchange(exc exchange.Exchange, base, quote string) bool {
	supported := false
	for _, pair := range exc.GetPairs() {
		if pair.Base == base && pair.Quote == quote {
			supported = true
			break
		}
	}
	return supported
}