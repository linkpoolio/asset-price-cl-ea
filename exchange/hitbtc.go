package exchange

import (
	"fmt"
)

type HitBtc struct {
	Exchange
	Pairs []*Pair
}

type HitBtcPair struct {
	Base  string `json:"baseCurrency"`
	Quote string `json:"quoteCurrency"`
}

type HitBtcTicker struct {
	VolumeQuote string `json:"volumeQuote"`
	Last  		string `json:"last"`
}

func (exc *HitBtc) GetResponse(base, quote string) (*Response, *Error) {
	var ticker HitBtcTicker
	config := exc.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/public/ticker/%s%s", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	return &Response{config.Name, ToFloat64(ticker.Last), ToFloat64(ticker.VolumeQuote)}, nil
}

func (exc *HitBtc) SetPairs() *Error {
	var pairs []HitBtcPair
	config := exc.GetConfig()
	err := HttpGet(config, "/public/symbol/", &pairs)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		exc.Pairs = append(exc.Pairs, &Pair{Base: pair.Base, Quote: pair.Quote})
	}
	return nil
}

func (exc *HitBtc) GetConfig() *Config {
	return &Config{
		Name: "HitBTC",
		BaseUrl: "https://api.hitbtc.com/api/2",
		Client: nil,
		Pairs: exc.Pairs}
}
