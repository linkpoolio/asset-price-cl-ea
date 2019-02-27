package exchange

import (
	"github.com/preichenberger/go-gdax"
	"fmt"
)

type Coinbase struct {
	Exchange
	Pairs []*Pair
}

func (exc *Coinbase) GetResponse(base, quote string) (*Response, *Error) {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*gdax.Client)

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, &Error{exc.GetConfig().Name, "500 ERROR", err.Error()}
	}

	return &Response{exc.GetConfig().Name, ticker.Price,  ticker.Volume * ticker.Price}, nil
}

func (exc *Coinbase) SetPairs() *Error {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*gdax.Client)

	products, err := client.GetProducts()
	if err != nil {
		return &Error{Exchange: exc.GetConfig().Name, Message: err.Error()}
	}
	for _, product := range products {
		exc.Pairs = append(exc.Pairs, &Pair{product.BaseCurrency, product.QuoteCurrency})
	}

	return nil
}

func (exc *Coinbase) GetConfig() *Config {
	return &Config{Name: "Coinbase", Client: gdax.NewClient("", "", ""), Pairs: exc.Pairs}
}