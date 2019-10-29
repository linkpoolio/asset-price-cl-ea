package exchange

import (
	"fmt"
	"strings"
)

type ZB struct {
	Exchange
}

type ZBTicker struct {
	Volume string `json:"vol"`
	Last   string `json:"last"`
}

type ZBMarket struct {
	Ticker ZBTicker `json:"ticker"`
}

func (exc *ZB) GetResponse(base, quote string) (*Response, error) {
	var market ZBMarket
	config := exc.GetConfig()
	err := exc.HttpGet(config, fmt.Sprintf("/ticker?market=%s_%s", base, quote), &market)
	if err != nil {
		return nil, err
	}
	volume := exc.toFloat64(market.Ticker.Volume) * exc.toFloat64(market.Ticker.Last)
	return &Response{Name: config.Name, Price: exc.toFloat64(market.Ticker.Last), Volume: volume}, nil
}

func (exc *ZB) RefreshPairs() error {
	var zbPairs map[string]interface{}
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/markets", &zbPairs)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for pair := range zbPairs {
		details := strings.Split(pair, "_")
		pairs = append(pairs, &Pair{Base: strings.ToUpper(details[0]), Quote: strings.ToUpper(details[1])})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *ZB) GetConfig() *Config {
	return &Config{Name: "ZB", BaseURL: "http://api.zb.com/data/v1"}
}
