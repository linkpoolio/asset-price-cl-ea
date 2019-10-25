package exchange

import (
	"fmt"
)

type Coinall struct {
	Exchange
	Pairs []*Pair
}

type CoinallPair struct {
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
}

type CoinallTicker struct {
	Last        string `json:"last"`
	QuoteVolume string `json:"quote_volume_24h"`
}

func (exc *Coinall) GetResponse(base, quote string) (*Response, *Error) {
	var ticker CoinallTicker
	config := exc.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/spot/v3/instruments/%s_%s/ticker", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	return &Response{config.Name, ToFloat64(ticker.Last), ToFloat64(ticker.QuoteVolume)}, nil
}

func (exc *Coinall) SetPairs() *Error {
	var pairs []CoinallPair
	config := exc.GetConfig()
	err := HttpGet(config, "/spot/v3/instruments", &pairs)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		exc.Pairs = append(exc.Pairs, &Pair{Base: pair.BaseCurrency, Quote: pair.QuoteCurrency})
	}
	return nil
}

func (exc *Coinall) GetConfig() *Config {
	return &Config{
		Name:    "Coinall",
		BaseUrl: "https://www.coinall.com/api",
		Client:  nil,
		Pairs:   exc.Pairs,
	}
}
