package exchange

import (
	"fmt"
	"github.com/preichenberger/go-gdax"
)

type Coinbase struct {
	Exchange
}

func (exc *Coinbase) GetResponse(base, quote string) (*Response, error) {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*gdax.Client)

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, &Error{exc.GetConfig().Name, "500 ERROR", err.Error()}
	}

	return &Response{exc.GetConfig().Name, ticker.Price, ticker.Volume * ticker.Price}, nil
}

func (exc *Coinbase) RefreshPairs() error {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*gdax.Client)

	products, err := client.GetProducts()
	if err != nil {
		return &Error{Exchange: exc.GetConfig().Name, Message: err.Error()}
	}

	var pairs []*Pair
	for _, product := range products {
		pairs = append(pairs, &Pair{product.BaseCurrency, product.QuoteCurrency})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Coinbase) GetConfig() *Config {
	return &Config{Name: "Coinbase", Client: gdax.NewClient("", "", "")}
}
