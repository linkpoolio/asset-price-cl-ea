package exchange

import (
	"fmt"
)

type Coinall struct {
	Exchange
}

type CoinallPair struct {
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
}

type CoinallTicker struct {
	Last        string `json:"last"`
	QuoteVolume string `json:"quote_volume_24h"`
}

func (exc *Coinall) GetResponse(base, quote string) (*Response, error) {
	var ticker CoinallTicker
	config := exc.GetConfig()
	err := exc.HttpGet(config, fmt.Sprintf("/spot/v3/instruments/%s_%s/ticker", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	return &Response{config.Name, exc.toFloat64(ticker.Last), exc.toFloat64(ticker.QuoteVolume)}, nil
}

func (exc *Coinall) RefreshPairs() error {
	var coinallPairs []CoinallPair
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/spot/v3/instruments", &coinallPairs)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, pair := range coinallPairs {
		pairs = append(pairs, &Pair{Base: pair.BaseCurrency, Quote: pair.QuoteCurrency})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Coinall) GetConfig() *Config {
	return &Config{Name: "Coinall", BaseURL: "https://www.coinall.com/api"}
}
