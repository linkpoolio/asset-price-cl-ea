package app

import (
	"errors"
	"fmt"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
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
	ID        string            `json:"id"`
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

func GetPrice(base, quote string) (*Output, error) {
	b := strings.ToUpper(base)
	q := strings.ToUpper(quote)

	output := Output{ID: fmt.Sprintf("%s-%s", b, q)}
	responses, ee := getExchangeResponses(b, q)
	if len(ee) > 0 {
		output.Warnings = ee
	} else if len(responses) == 0 {
		return nil, errors.New("No exchanges support that trading pair")
	}

	p, v := aggregateResponses(responses)
	output.Price = formatFloat(p)
	output.Volume = formatFloat(v)

	if quote == "USD" {
		output.USDPrice = null.StringFrom(formatFloat(p))
	} else {
		qup, ee := getQuoteUSDPrice(q)
		output.Warnings = append(output.Warnings, ee...)
		if qup != 0 {
			output.USDPrice = null.StringFrom(formatFloat(qup*p))
		}
	}

	for _, response := range responses {
		if response != nil {
			output.Exchanges = append(output.Exchanges, response.Name)
		}
	}

	return &output, nil
}

func StartPairsTicker(c *Config) {
	setExchangePairs()

	ticker := time.NewTicker(c.TickerInterval)
	go func() {
		for range ticker.C {
			setExchangePairs()
			log.Print("Trading pairs refreshed")
		}
	}()
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
				Message:  fmt.Sprintf("No exchange supports the %s-USD pair for fetching usd price", quote),
				Status:   "400",
			},
		}
	}
	p, _ := aggregateResponses(responses)
	return p, err
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
				}).Error("Error from exchange on setting pairs")
			}
		}(exc)
	}

	wg.Wait()
}

func formatFloat(float float64) string {
	str := strconv.FormatFloat(float, 'f', -1, 64)
	if str == "NaN" {
		return "0"
	}
	return str
}