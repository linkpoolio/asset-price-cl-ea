package exchange

import (
	"fmt"
	"log"
	"strings"
)

type Bitfinex struct {
	Exchange
}

type BitfinexTicker struct {
	Volume 	  string `json:"volume"`
	LastPrice string `json:"last_price"`
}

var bitfinexPairs []*Pair

func (exchange *Bitfinex) GetResponse(base, quote string) (*Response, *Error) {
	var ticker BitfinexTicker
	config := exchange.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/pubticker/%s%s", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	volume := ToFloat64(ticker.Volume) * ToFloat64(ticker.LastPrice)
	return &Response{Name: config.Name, Price: ToFloat64(ticker.LastPrice), Volume: volume}, nil
}

func (exchange *Bitfinex) SetPairs() {
	var pairs []string
	config := exchange.GetConfig()
	err := HttpGet(config, "/symbols", &pairs)
	if err != nil {
		log.Fatal(err)
	}
	// We have to assume all BTC pairs are 3char base, 3char quote. No base/quote given in API.
	for _, pair := range pairs {
		if len(pair) == 6 {
			bitfinexPairs = append(
				bitfinexPairs,
				&Pair{Base: strings.ToUpper(pair[0:3]), Quote: strings.ToUpper(pair[3:6])})
		}
	}
}

func (exchange *Bitfinex) GetConfig() *Config {
	return &Config{
		Name: "Bitfinex",
		BaseUrl: "https://api.bitfinex.com/v1",
		Client: nil,
		Pairs: bitfinexPairs}
}
