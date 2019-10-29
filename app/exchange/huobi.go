package exchange

import (
	"fmt"
	"strings"
)

type Huobi struct {
	Exchange
}

type HuobiPair struct {
	Base  string `json:"base-currency"`
	Quote string `json:"quote-currency"`
}

type HuobiPairs struct {
	Data []HuobiPair `json:"data"`
}

type HuobiTicker struct {
	Volume float64 `json:"vol"`
	Close  float64 `json:"close"`
}

type HuobiMarket struct {
	Ticker HuobiTicker `json:"tick"`
}

func (exc *Huobi) GetResponse(base, quote string) (*Response, error) {
	var market HuobiMarket
	config := exc.GetConfig()
	err := exc.HttpGet(
		config,
		fmt.Sprintf(
			"/market/detail/merged?symbol=%s%s",
			strings.ToLower(base),
			strings.ToLower(quote)),
		&market)
	if err != nil {
		return nil, err
	}
	return &Response{config.Name, market.Ticker.Close, market.Ticker.Volume}, nil
}

func (exc *Huobi) RefreshPairs() error {
	var huobiPairs HuobiPairs
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/v1/common/symbols", &huobiPairs)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, pair := range huobiPairs.Data {
		pairs = append(pairs, &Pair{Base: strings.ToUpper(pair.Base), Quote: strings.ToUpper(pair.Quote)})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Huobi) GetConfig() *Config {
	return &Config{Name: "Huobi", BaseURL: "https://api.huobi.pro"}
}
