package exchange

import (
	"fmt"
	"strings"
)

type ZB struct {
	Exchange
	Pairs []*Pair
}

type ZBTicker struct {
	Volume string `json:"vol"`
	Last   string `json:"last"`
}

type ZBMarket struct {
	Ticker ZBTicker `json:"ticker"`
}

func (exc *ZB) GetResponse(base, quote string) (*Response, *Error) {
	var market ZBMarket
	config := exc.GetConfig()
	err := HttpGet(config, fmt.Sprintf("/ticker?market=%s_%s", base, quote), &market)
	if err != nil {
		return nil, err
	}
	volume := ToFloat64(market.Ticker.Volume) * ToFloat64(market.Ticker.Last)
	return &Response{Name: config.Name, Price: ToFloat64(market.Ticker.Last), Volume: volume}, nil
}

func (exc *ZB) SetPairs() *Error {
	var pairs map[string]interface{}
	config := exc.GetConfig()
	err := HttpGet(config, "/markets", &pairs)
	if err != nil {
		return err
	}

	for pair := range pairs {
		details := strings.Split(pair, "_")
		exc.Pairs = append(exc.Pairs, &Pair{Base: strings.ToUpper(details[0]), Quote: strings.ToUpper(details[1])})
	}
	return nil
}

func (exc *ZB) GetConfig() *Config {
	return &Config{
		Name:    "ZB",
		BaseUrl: "http://api.zb.com/data/v1",
		Client:  nil,
		Pairs:   exc.Pairs}
}
