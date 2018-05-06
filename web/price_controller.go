package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
	"fmt"
	"sync"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/guregu/null.v3"
	"strconv"
)

func GetResponse(w rest.ResponseWriter, r *rest.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	// Unmarshal to CL's RunResult type
	var runResult RunResult
	err := json.Unmarshal(bytes, &runResult)
	if err != nil {
		writeErrorResult(w, http.StatusInternalServerError, &runResult, err)
		return
	}

	// Call the exchanges concurrently
	responses, errors := getExchangeResponses(runResult.Params.Base, runResult.Params.Quote)
	if len(errors) > 0 {
		writeErrorResult(
			w,
			http.StatusInternalServerError,
			&runResult,
			fmt.Errorf("errors given when getting prices from exchanges"))
		return
	} else if len(responses) == 0 {
		writeErrorResult(
			w,
			http.StatusBadRequest,
			&runResult,
			fmt.Errorf("no exchanges support that trading pair"))
		return
	}

	// Calculate the weighted average based on volume
	var price float64
	var volume float64

	output := Output{Id: fmt.Sprintf("%s-%s", runResult.Params.Base, runResult.Params.Quote)}
	for _, response := range responses {
		output.Exchanges = append(output.Exchanges, response.Name)
		volume += response.Volume
	}
	for _, response := range responses {
		price += (response.Volume / volume) * response.Price
	}
	output.Price = strconv.FormatFloat(price, 'f', -1, 64)
	output.Volume = strconv.FormatFloat(volume, 'f', -1, 64)
	params := Params{ runResult.Params.Input, output }
	runResult.Params = params

	w.WriteJson(runResult)
}

func SetTradingPairs() {
	exchanges := exchange.GetSupportedExchanges()

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	for _, exc := range exchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			exc.SetPairs()
		}(exc)
	}

	wg.Wait()
}

func writeErrorResult(w rest.ResponseWriter, statusCode int, rr *RunResult, err error) {
	errorMessage := null.String{}
	errorMessage.String = err.Error()
	errorMessage.Valid = true

	rr.Pending = false
	rr.ErrorMessage = errorMessage

	w.WriteHeader(statusCode)
	w.WriteJson(rr)
}

func getExchangeResponses(base, quote string) ([]*exchange.Response, []*exchange.Error) {
	exchanges := getExchangesWithPairSupport(base, quote)

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	var responses []*exchange.Response
	var errors []*exchange.Error

	for _, exc := range exchanges {
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

func getExchangesWithPairSupport(base, quote string) []exchange.Exchange {
	exchanges := exchange.GetSupportedExchanges()

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	var supported []exchange.Exchange
	for _, exc := range exchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			for _, pair := range exc.GetConfig().Pairs {
				if pair.Base == base && pair.Quote == quote {
					supported = append(supported, exc)
					break
				}
			}
		}(exc)
	}

	wg.Wait()

	return supported
}