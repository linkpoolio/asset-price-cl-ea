package exchange

import (
	"fmt"
	"log"
)

type HitBtc struct {
	Exchange
}

type HitBtcPair struct {
	Base  string `json:"baseCurrency"`
	Quote string `json:"quoteCurrency"`
}

type HitBtcTicker struct {
	VolumeQuote string `json:"volumeQuote"`
	Last  		string `json:"last"`
}

var hitBtcPairs []*Pair

func (exchange *HitBtc) GetResponse(base, quote string) (*Response, *Error) {
	var ticker HitBtcTicker
	config := exchange.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/public/ticker/%s%s", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	return &Response{config.Name, ToFloat64(ticker.Last), ToFloat64(ticker.VolumeQuote)}, nil
}

func (exchange *HitBtc) SetPairs() {
	var pairs []HitBtcPair
	config := exchange.GetConfig()
	err := HttpGet(config, "/public/symbol/", &pairs)
	if err != nil {
		log.Fatal(err)
	}
	for _, pair := range pairs {
		hitBtcPairs = append(hitBtcPairs, &Pair{Base: pair.Base, Quote: pair.Quote})
	}
}

func (exchange *HitBtc) GetConfig() *Config {
	return &Config{
		Name: "HitBTC",
		BaseUrl: "https://api.hitbtc.com/api/2",
		Client: nil,
		Pairs: hitBtcPairs}
}
