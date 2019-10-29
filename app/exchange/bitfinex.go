package exchange

import (
	"fmt"
	"strings"
)

type Bitfinex struct {
	Exchange
}

type BitfinexTicker struct {
	Volume    string `json:"volume"`
	LastPrice string `json:"last_price"`
}

func (exc *Bitfinex) GetResponse(base, quote string) (*Response, error) {
	var ticker BitfinexTicker
	config := exc.GetConfig()
	err := exc.HttpGet(config, fmt.Sprintf("/pubticker/%s%s", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	volume := exc.toFloat64(ticker.Volume) * exc.toFloat64(ticker.LastPrice)
	return &Response{Name: config.Name, Price: exc.toFloat64(ticker.LastPrice), Volume: volume}, nil
}

func (exc *Bitfinex) RefreshPairs() error {
	var symbols []string
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/symbols", &symbols)
	if err != nil {
		return err
	}

	var pairs []*Pair
	// We have to assume all BTC symbols are 3char base, 3char quote. No base/quote given in API.
	for _, symbol := range symbols {
		if len(symbol) == 6 {
			pairs = append(
				pairs,
				&Pair{Base: strings.ToUpper(symbol[0:3]), Quote: strings.ToUpper(symbol[3:6])})
		}
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Bitfinex) GetConfig() *Config {
	return &Config{
		Name:    "Bitfinex",
		BaseURL: "https://api.bitfinex.com/v1",
		Client:  nil,
	}
}
