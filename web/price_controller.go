package web

import (
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Log struct {
	Message string
}

type Input struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type Output struct {
	Id        string            `json:"id"`
	Price     string            `json:"price"`
	Volume    string            `json:"volume"`
	USDPrice  null.String       `json:"usdPrice"`
	Exchanges []string          `json:"exchanges"`
	Warnings  []*exchange.Error `json:"warnings"`
}

type Params struct {
	Input
	Output
}

type RunResult struct {
	JobRunID     string      `json:"jobRunId"`
	Params       Params      `json:"data"`
	Status       string      `json:"status"`
	ErrorMessage null.String `json:"error"`
	Pending      bool        `json:"pending"`
}

func GetResponse(w rest.ResponseWriter, r *rest.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	var runResult RunResult
	err := json.Unmarshal(bytes, &runResult)
	if err != nil {
		writeErrorResult(w, http.StatusInternalServerError, &runResult, err)
		return
	}

	b := strings.ToUpper(runResult.Params.Base)
	q := strings.ToUpper(runResult.Params.Quote)

	output := Output{Id: fmt.Sprintf("%s-%s", runResult.Params.Base, runResult.Params.Quote)}
	responses, ee := getExchangeResponses(b, q)
	if len(ee) > 0 {
		output.Warnings = append(runResult.Params.Warnings, ee...)
	} else if len(responses) == 0 {
		writeErrorResult(
			w,
			http.StatusBadRequest,
			&runResult,
			fmt.Errorf("no exchanges support that trading pair"))
		return
	}

	p, v := aggregateResponses(responses)
	output.Price = strconv.FormatFloat(p, 'f', -1, 64)
	output.Volume = strconv.FormatFloat(v, 'f', -1, 64)

	qup, ee := getQuoteUSDPrice(q)
	if strings.Contains(q, "USD") {
		output.USDPrice = null.StringFrom(output.Price)
	} else if len(ee) == 0 {
		output.USDPrice = null.StringFrom(strconv.FormatFloat(qup*p, 'f', -1, 64))
	} else {
		output.Warnings = append(runResult.Params.Warnings, ee...)
	}

	for _, response := range responses {
		if response != nil {
			output.Exchanges = append(output.Exchanges, response.Name)
		}
	}

	params := Params{runResult.Params.Input, output}
	runResult.Params = params

	w.WriteJson(runResult)
}

func StartPairsTicker() {
	setExchangePairs()

	ticker := time.NewTicker(Config.TickerInterval)
	go func() {
		for range ticker.C {
			setExchangePairs()
			log.Print("trading pairs refreshed")
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

func getQuoteUSDPrice(quote string) (float64, []*exchange.Error) {
	responses, err := getExchangeResponses(quote, "USD")
	if len(responses) == 0 {
		return 0, []*exchange.Error{
			{
				Exchange: "N/A",
				Message:  fmt.Sprintf("no exchange supports the %s-USD pair for fetching usd price", quote),
				Status:   "400",
			},
		}
	}
	if err != nil {
		return 0, err
	}
	p, _ := aggregateResponses(responses)
	return p, nil
}

func aggregateResponses(responses []*exchange.Response) (float64, float64) {
	var price float64
	var volume float64
	for _, response := range responses {
		if response != nil {
			volume += response.Volume
		}
	}
	for _, response := range responses {
		if response != nil {
			price += (response.Volume / volume) * response.Price
		}
	}
	return price, volume
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
			err := exc.SetPairs()
			if err != nil {
				log.WithFields(log.Fields{
					"exchange": err.Exchange,
					"msg":      err.Message,
				}).Error("error from exchange on setting pairs")
			}
		}(exc)
	}

	wg.Wait()
}
