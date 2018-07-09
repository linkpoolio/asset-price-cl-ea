package exchange

import (
	"fmt"
	"log"
	"strings"
)

type ZB struct {
	Exchange
}

type ZBTicker struct {
	Volume 	string `json:"vol"`
	Last 	string `json:"last"`
}

type ZBMarket struct {
	Ticker ZBTicker `json:"ticker"`
}

var zbPairs []*Pair

func (exchange *ZB) GetResponse(base, quote string) (*Response, *Error) {
	var market ZBMarket
	config := exchange.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/ticker?market=%s_%s", base, quote), &market)
	if err != nil {
		return nil, err
	}
	volume := ToFloat64(market.Ticker.Volume) * ToFloat64(market.Ticker.Last)
	return &Response{Name: config.Name, Price: ToFloat64(market.Ticker.Last), Volume: volume}, nil
}

func (exchange *ZB) SetPairs() {
	var pairs map[string]interface{}
	config := exchange.GetConfig()
	err := HttpGet(config, "/markets", &pairs)
	if err != nil {
		log.Fatal(err)
	}

	for pair := range pairs {
		details := strings.Split(pair, "_")
		zbPairs = append(zbPairs, &Pair{Base: strings.ToUpper(details[0]), Quote: strings.ToUpper(details[1])})
	}
}

func (exchange *ZB) GetConfig() *Config {
	return &Config{
		Name: "ZB",
		BaseUrl: "http://api.zb.com/data/v1",
		Client: nil,
		Pairs: zbPairs}
}
