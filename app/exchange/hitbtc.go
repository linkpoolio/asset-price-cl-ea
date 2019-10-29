package exchange

import (
	"fmt"
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
	Last        string `json:"last"`
}

func (exc *HitBtc) GetResponse(base, quote string) (*Response, error) {
	var ticker HitBtcTicker
	config := exc.GetConfig()
	err := exc.HttpGet(config, fmt.Sprintf("/public/ticker/%s%s", base, quote), &ticker)
	if err != nil {
		return nil, err
	}
	return &Response{config.Name, exc.toFloat64(ticker.Last), exc.toFloat64(ticker.VolumeQuote)}, nil
}

func (exc *HitBtc) RefreshPairs() error {
	var hitBtcPairs []HitBtcPair
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/public/symbol/", &hitBtcPairs)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, pair := range hitBtcPairs {
		pairs = append(pairs, &Pair{Base: pair.Base, Quote: pair.Quote})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *HitBtc) GetConfig() *Config {
	return &Config{Name: "HitBTC", BaseURL: "https://api.hitbtc.com/api/2"}
}
