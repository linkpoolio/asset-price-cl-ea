package exchange

import (
	"fmt"
	"strings"
)

type Kraken struct {
	Exchange
}

type KrakenPairs struct {
	Result map[string]*KrakenPair `json:"result"`
}

type KrakenPair struct {
	WSName string `json:"wsname"`
}

type KrakenTickers struct {
	Result map[string]*KrakenTicker `json:"result"`
}

type KrakenTicker struct {
	LastTrade []string `json:"c"`
	Volume    []string `json:"v"`
}

func (exc *Kraken) GetResponse(base, quote string) (*Response, error) {
	var kt KrakenTickers
	config := exc.GetConfig()
	excErr := exc.HttpGet(config, fmt.Sprintf("/Ticker?pair=%s%s", base, quote), &kt)
	if excErr != nil {
		return nil, excErr
	} else if len(kt.Result) == 0 {
		return nil, &Error{Exchange: config.Name, Status: "400", Message: "Ticker call returned invalid data"}
	}

	for p, t := range kt.Result {
		if strings.Contains(p, base) && strings.Contains(p, quote) {
			return &Response{
				Name:   config.Name,
				Price:  exc.toFloat64(t.LastTrade[0]),
				Volume: exc.toFloat64(t.Volume[1]),
			}, nil
		}
	}

	return nil, &Error{Exchange: config.Name, Status: "400", Message: "Ticker call didn't return the requested pair"}
}

func (exc *Kraken) RefreshPairs() error {
	var kp KrakenPairs
	config := exc.GetConfig()

	err := exc.HttpGet(config, "/AssetPairs", &kp)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, p := range kp.Result {
		pSlice := strings.Split(p.WSName, "/")
		if len(pSlice) == 2 {
			pairs = append(pairs, &Pair{pSlice[0], pSlice[1]})
		}
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Kraken) GetConfig() *Config {
	return &Config{Name: "Kraken", BaseURL: "https://api.kraken.com/0/public"}
}
