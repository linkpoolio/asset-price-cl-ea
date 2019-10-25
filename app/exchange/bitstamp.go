package exchange

import (
	"fmt"
	"strings"
)

type Bitstamp struct {
	Exchange
	Pairs []*Pair
}

type BitstampModel struct {
	Last   string `json:"last"`
	Volume string `json:"volume"`
}

type BitstampPair struct {
	Name    string `json:"name"`
	Trading string `json:"trading"`
}

func (exc *Bitstamp) GetResponse(base, quote string) (*Response, *Error) {
	var bst BitstampModel
	config := exc.GetConfig()
	excErr := HttpGet(config, fmt.Sprintf("/ticker/%s%s", base, quote), &bst)
	if excErr != nil {
		return nil, excErr
	}
	volume := ToFloat64(bst.Volume) * ToFloat64(bst.Last)
	return &Response{exc.GetConfig().Name, ToFloat64(bst.Last), volume}, nil
}

func (exc *Bitstamp) SetPairs() *Error {
	var pairs []BitstampPair
	config := exc.GetConfig()
	err := HttpGet(config, "/trading-pairs-info/", &pairs)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		if pair.Trading == "Enabled" {
			currencies := strings.Split(pair.Name, "/")
			exc.Pairs = append(exc.Pairs, &Pair{currencies[0], currencies[1]})
		}
	}
	return nil
}

func (exc *Bitstamp) GetConfig() *Config {
	return &Config{BaseUrl: "https://bitstamp.net/api/v2", Name: "Bitstamp", Pairs: exc.Pairs}
}
