package exchange

import (
	"github.com/preichenberger/go-gdax"
	"fmt"
)

type GDAX struct {
	Exchange
}

func (exchange GDAX) GetResponse(base, quote string) (*Response, *Error) {
	client := gdax.NewClient("", "", "")

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, "500 ERROR", err.Error()}
	}

	return &Response{exchange.GetConfig().Name, ticker.Price,  ticker.Volume}, nil
}

func (exchange GDAX) GetConfig() *Config {
	return &Config{Name: "GDAX", SupportedBases: []string{"BTC", "BCH", "LTC"}, SupportedQuotes: []string{"USD", "EUR"}}
}