package exchange

import (
	"fmt"
	"strings"
)

type COSS struct {
	Exchange
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

func (exc *COSS) GetResponse(base, quote string) (*Response, error) {
	var ticker COSSTicker
	config := exc.GetConfig()
	config.BaseURL = "https://exchange.coss.io/api" // Different endpoint between their engine/exchange API
	err := exc.HttpGet(config, fmt.Sprintf("/getmarketsummaries"), &ticker)
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
		Status:   "400",
		Message:  "pair given wasn't found within market summaries",
	}
}

func (exc *COSS) RefreshPairs() error {
	var cossPairs COSSPairs
	config := exc.GetConfig()
	err := exc.HttpGet(config, "/exchange-info", &cossPairs)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, pair := range cossPairs.Symbols {
		pairArr := strings.Split(pair.Symbol, "_")
		if len(pairArr) == 2 {
			pairs = append(pairs, &Pair{Base: pairArr[0], Quote: pairArr[1]})
		}
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *COSS) GetConfig() *Config {
	return &Config{Name: "COSS", BaseURL: "https://trade.coss.io/c/api/v1"}
}
