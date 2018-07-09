package exchange

import (
	"fmt"
	"strings"
	"log"
)

type Bitstamp struct {
	Exchange
}

type BitstampModel struct {
	Last string `json:"last"`
	Volume string `json:"volume"`
}

type BitstampPair struct {
	Name string `json:"name"`
	Trading string `json:"trading"`
}

var bitstampPairs []*Pair

func (exchange *Bitstamp) GetResponse(base, quote string) (*Response, *Error) {
	var bst BitstampModel
	config := exchange.GetConfig()
	excErr := HttpGet(config, fmt.Sprintf("/ticker/%s%s", base, quote), &bst)
	if excErr != nil {
		return nil, excErr
	}
	volume := ToFloat64(bst.Volume) * ToFloat64(bst.Last)
	return &Response{exchange.GetConfig().Name, ToFloat64(bst.Last), volume}, nil
}

func (exchange *Bitstamp) SetPairs() {
	var pairs []BitstampPair
	config := exchange.GetConfig()
	err := HttpGet(config, "/trading-pairs-info/", &pairs)
	if err != nil {
		log.Fatal(err)
	}
	for _, pair := range pairs {
		if pair.Trading == "Enabled" {
			currencies := strings.Split(pair.Name, "/")
			bitstampPairs = append(bitstampPairs, &Pair{currencies[0], currencies[1]})
		}
	}
}

func (exchange *Bitstamp) GetConfig() *Config {
	return &Config{BaseUrl: "https://bitstamp.net/api/v2", Name: "Bitstamp", Pairs: bitstampPairs}
}