package exchange

import (
	"fmt"
)

type Bittrex struct {
	Exchange
}

type BittrexTicker struct {
	Volume    string `json:"volume"`
	LastPrice string `json:"last_price"`
}

type BittrexMarket struct {
	MarketCurrency string `json:"MarketCurrency"`
	BaseCurrency   string `json:"BaseCurrency"`
}

type BittrexMarkets struct {
	Result []BittrexMarket `json:"result"`
}

type BittrexSummary struct {
	Last   float64 `json:"Last"`
	Volume float64 `json:"Volume"`
}

type BittrexSummaries struct {
	Result []BittrexSummary `json:"result"`
}

func (exc *Bittrex) GetResponse(base, quote string) (*Response, error) {
	var summaries BittrexSummaries
	config := exc.GetConfig()
	err := exc.HttpGet(config, fmt.Sprintf("/public/getmarketsummary?market=%s-%s", base, quote), &summaries)
	if err != nil {
		return nil, err
	}
	return &Response{Name: config.Name, Price: summaries.Result[0].Last, Volume: summaries.Result[0].Volume}, nil
}

func (exc *Bittrex) RefreshPairs() error {
	var markets BittrexMarkets
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/public/getmarkets", &markets)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, pair := range markets.Result {
		pairs = append(pairs, &Pair{Base: pair.BaseCurrency, Quote: pair.MarketCurrency})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Bittrex) GetConfig() *Config {
	return &Config{Name: "Bittrex", BaseURL: "https://bittrex.com/api/v1.1"}
}
