package exchange

import (
	"github.com/preichenberger/go-gdax"
	"fmt"
	"log"
)

type GDAX struct {
	Exchange
}

var gdaxPairs []*Pair

func (exchange GDAX) GetResponse(base, quote string) (*Response, *Error) {
	clientInterface := exchange.GetConfig().Client
	client := clientInterface.(*gdax.Client)

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, "500 ERROR", err.Error()}
	}

	return &Response{exchange.GetConfig().Name, ticker.Price,  ticker.Volume}, nil
}

func (exchange GDAX) SetPairs() {
	clientInterface := exchange.GetConfig().Client
	client := clientInterface.(*gdax.Client)

	products, err := client.GetProducts()
	if err != nil {
		log.Fatal(err)
	}
	for _, product := range products {
		gdaxPairs = append(gdaxPairs, &Pair{product.BaseCurrency, product.QuoteCurrency})
	}
}

func (exchange GDAX) GetConfig() *Config {
	return &Config{Name: "GDAX", Client: gdax.NewClient("", "", ""), Pairs: gdaxPairs}
}