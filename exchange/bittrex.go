package exchange

import (
	"fmt"
	"log"
)

type Bittrex struct {
	Exchange
	Pairs []*Pair
}

type BittrexTicker struct {
	Volume 	  string `json:"volume"`
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

func (exc *Bittrex) GetResponse(base, quote string) (*Response, *Error) {
	var summaries BittrexSummaries
	config := exc.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/public/getmarketsummary?market=%s-%s", base, quote), &summaries)
	if err != nil {
		return nil, err
	}
	return &Response{Name: config.Name, Price: summaries.Result[0].Last, Volume: summaries.Result[0].Volume}, nil
}

func (exc *Bittrex) SetPairs() {
	var markets BittrexMarkets
	config := exc.GetConfig()
	err := HttpGet(config, "/public/getmarkets", &markets)
	if err != nil {
		log.Fatal(err)
	}
	for _, pair := range markets.Result {
		exc.Pairs = append(exc.Pairs, &Pair{Base: pair.BaseCurrency, Quote: pair.MarketCurrency})
	}
}

func (exc *Bittrex) GetConfig() *Config {
	return &Config{
		Name: "Bittrex",
		BaseUrl: "https://bittrex.com/api/v1.1",
		Client: nil,
		Pairs: exc.Pairs}
}
