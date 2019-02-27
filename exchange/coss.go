package exchange

import (
	"fmt"
	"strings"
)

type COSS struct {
	Exchange
	Pairs []*Pair
}

type COSSPairs struct {
	Symbols []COSSPair `json:"symbols"`
}

type COSSPair struct {
	Symbol string `json:"symbol"`
}

type COSSTicker struct {
	Results []*COSSResult `json:"result"`
}

type COSSResult struct {
	MarketName string  `json:"MarketName"`
	Last       float64 `json:"Last"`
	BaseVolume float64 `json:"BaseVolume"`
}

func (exc *COSS) GetResponse(base, quote string) (*Response, *Error) {
	var ticker COSSTicker
	config := exc.GetConfig()
	config.BaseUrl = "https://exchange.coss.io/api" // Different endpoint between their engine/exchange API
	err := HttpGet(config, fmt.Sprintf("/getmarketsummaries"), &ticker)
	if err != nil {
		return nil, err
	}
	mn := fmt.Sprintf("%s-%s", base, quote)
	for _, result := range ticker.Results {
		if result.MarketName == mn {
			return &Response{config.Name, result.Last, result.BaseVolume}, nil
		}
	}
	return nil, &Error{
		Exchange: config.Name,
		Status: "400",
		Message: "pair given wasn't found within market summaries",
	}
}

func (exc *COSS) SetPairs() *Error {
	var pairs COSSPairs
	config := exc.GetConfig()
	err := HttpGet(config, "/exchange-info", &pairs)
	if err != nil {
		return err
	}
	for _, pair := range pairs.Symbols {
		pairArr := strings.Split(pair.Symbol, "_")
		if len(pairArr) == 2 {
			exc.Pairs = append(exc.Pairs, &Pair{Base: pairArr[0], Quote: pairArr[1]})
		}
	}
	return nil
}

func (exc *COSS) GetConfig() *Config {
	return &Config{
		Name:    "COSS",
		BaseUrl: "https://trade.coss.io/c/api/v1",
		Client:  nil,
		Pairs:   exc.Pairs,
	}
}
