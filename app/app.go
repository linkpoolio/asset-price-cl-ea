package app

import (
	"errors"
	"fmt"
	"github.com/linkpoolio/asset-price-cl-ea/app/exchange"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Output struct {
	ID        string            `json:"id"`
	Result    float64           `json:"result"`
	Price     string            `json:"price"`
	Volume    string            `json:"volume"`
	USDPrice  null.String       `json:"usdPrice"`
	Exchanges []string          `json:"exchanges"`
	Warnings  []error           `json:"warnings,omitempty"`
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
	output.Result = p
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
	if c == nil {
		return
	}

	ticker := time.NewTicker(c.TickerInterval)
	go func() {
		for range ticker.C {
			setExchangePairs()
			log.Print("Trading pairs refreshed")
		}
	}()
}

func getExchangeResponses(base, quote string) ([]*exchange.Response, []error) {
	exchanges := getExchangesWithPairSupport(base, quote)

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	var responses []*exchange.Response
	var errs []error

	mutex := sync.Mutex{}
	for _, exc := range exchanges {
		go func(exc exchange.Interface) {
			defer wg.Done()
			response, err := exc.GetResponse(base, quote)
			mutex.Lock()
			if err != nil {
				errs = append(errs, err)
			}
			responses = append(responses, response)
			mutex.Unlock()
		}(exc)
	}

	wg.Wait()

	return responses, errs
}

func getQuoteUSDPrice(quote string) (float64, []error) {
	responses, errs := getExchangeResponses(quote, "USD")
	if len(responses) == 0 {
		errs = append(errs, &exchange.Error{
			Exchange: "N/A",
			Message:  fmt.Sprintf("No exchange supports the %s-USD pair for fetching usd price", quote),
			Status:   "400",
		})
	}
	p, _ := aggregateResponses(responses)
	return p, errs
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

func getExchangesWithPairSupport(base, quote string) []exchange.Interface {
	exchanges := exchange.GetSupportedExchanges()

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	mutex := sync.Mutex{}
	var supported []exchange.Interface
	for _, exc := range exchanges {
		go func(exc exchange.Interface) {
			defer wg.Done()
			for _, pair := range exc.GetPairs() {
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
		go func(exc exchange.Interface) {
			defer wg.Done()
			err := exc.RefreshPairs()
			if err != nil {
				log.Error(err)
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