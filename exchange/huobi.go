package exchange

import (
	"fmt"
	"strings"
)

type Huobi struct {
	Exchange
	Pairs []*Pair
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

func (exc *Huobi) GetResponse(base, quote string) (*Response, *Error) {
	var market HuobiMarket
	config := exc.GetConfig()
	err := HttpGet(
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

func (exc *Huobi) SetPairs() *Error {
	var pairs HuobiPairs
	config := exc.GetConfig()
	err := HttpGet(config, "/v1/common/symbols", &pairs)
	if err != nil {
		return err
	}
	for _, pair := range pairs.Data {
		exc.Pairs = append(exc.Pairs, &Pair{Base: strings.ToUpper(pair.Base), Quote: strings.ToUpper(pair.Quote)})
	}
	return nil
}

func (exc *Huobi) GetConfig() *Config {
	return &Config{
		Name: "Huobi",
		BaseUrl: "https://api.huobi.pro",
		Client: nil,
		Pairs: exc.Pairs}
}
