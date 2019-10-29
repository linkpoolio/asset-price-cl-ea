package exchange

import (
	"fmt"
	"strings"
)

type Bitstamp struct {
	Exchange
}

type BitstampModel struct {
	Last   string `json:"last"`
	Volume string `json:"volume"`
}

type BitstampPair struct {
	Name    string `json:"name"`
	Trading string `json:"trading"`
}

func (exc *Bitstamp) GetResponse(base, quote string) (*Response, error) {
	var bst BitstampModel
	config := exc.GetConfig()
	excErr := exc.HttpGet(config, fmt.Sprintf("/ticker/%s%s", base, quote), &bst)
	if excErr != nil {
		return nil, excErr
	}
	volume := exc.toFloat64(bst.Volume) * exc.toFloat64(bst.Last)
	return &Response{exc.GetConfig().Name, exc.toFloat64(bst.Last), volume}, nil
}

func (exc *Bitstamp) RefreshPairs() error {
	var bitstampPairs []BitstampPair
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/trading-bitstampPairs-info/", &bitstampPairs)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, pair := range bitstampPairs {
		if pair.Trading == "Enabled" {
			currencies := strings.Split(pair.Name, "/")
			pairs = append(pairs, &Pair{currencies[0], currencies[1]})
		}
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Bitstamp) GetConfig() *Config {
	return &Config{BaseURL: "https://bitstamp.net/api/v2", Name: "Bitstamp"}
}
