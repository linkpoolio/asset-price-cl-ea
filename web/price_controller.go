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
	"strings"
	"time"
	"log"
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

func GetResponse(w rest.ResponseWriter, r *rest.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	var runResult RunResult
	err := json.Unmarshal(bytes, &runResult)
	if err != nil {
		writeErrorResult(w, http.StatusInternalServerError, &runResult, err)
		return
	}

	output := Output{Id: fmt.Sprintf("%s-%s", runResult.Params.Base, runResult.Params.Quote)}
	responses, errors := getExchangeResponses(
		strings.ToUpper(runResult.Params.Base),
		strings.ToUpper(runResult.Params.Quote))
	if len(errors) > 0 {
		output.Errors = append(runResult.Params.Errors, errors...)
	} else if len(responses) == 0 {
		writeErrorResult(
			w,
			http.StatusBadRequest,
			&runResult,
			fmt.Errorf("no exchanges support that trading pair"))
		return
	}

	var price float64
	var volume float64

	for _, response := range responses {
		if response != nil {
			output.Exchanges = append(output.Exchanges, response.Name)
			volume += response.Volume
		}
	}
	for _, response := range responses {
		if response != nil {
			price += (response.Volume / volume) * response.Price
		}
	}
	output.Price = strconv.FormatFloat(price, 'f', -1, 64)
	output.Volume = strconv.FormatFloat(volume, 'f', -1, 64)
	params := Params{ runResult.Params.Input, output }
	runResult.Params = params

	w.WriteJson(runResult)
}

func StartPairsTicker() {
	setExchangePairs()

	ticker := time.NewTicker(Config.TickerInterval)
	go func() {
		for range ticker.C {
			setExchangePairs()
			log.Print("Trading pairs refreshed")
		}
	}()
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

	mutex := sync.Mutex{}
	for _, exc := range exchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			response, err := exc.GetResponse(base, quote)
			mutex.Lock()
			if err != nil {
				errors = append(errors, err)
			}
			responses = append(responses, response)
			mutex.Unlock()
		}(exc)
	}

	wg.Wait()

	return responses, errors
}

func getExchangesWithPairSupport(base, quote string) []exchange.Exchange {
	exchanges := exchange.GetSupportedExchanges()

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	mutex := sync.Mutex{}
	var supported []exchange.Exchange
	for _, exc := range exchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			for _, pair := range exc.GetConfig().Pairs {
				if pair.Base == base && pair.Quote == quote {
					mutex.Lock()
					supported = append(supported, exc)
					mutex.Unlock()
					break
				}
			}
		}(exc)
	}

	wg.Wait()

	return supported
}

func setExchangePairs() {
	var wg sync.WaitGroup

	exchanges := exchange.GetSupportedExchanges()
	wg.Add(len(exchanges))

	for _, exc := range exchanges {
		go func(exc exchange.Exchange) {
			defer wg.Done()
			exc.SetPairs()
		}(exc)
	}

	wg.Wait()
}