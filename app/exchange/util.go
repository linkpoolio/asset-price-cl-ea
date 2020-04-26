package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	exchanges = []Interface{
		&Binance{},
		&Bitfinex{},
		//&Bitstamp{}, (Extreme rate limit)
		&Bittrex{},
		&Coinall{},
		&Coinbase{},
		&Gemini{},
		&Huobi{},
		&HitBtc{},
		&Kraken{},
		&ZB{},
	}
)

type Config struct {
	Name    string
	BaseURL string
	Client  interface{}
}

type Error struct {
	Exchange string
	Status   string
	Message  string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s (status: %s)", e.Exchange, e.Message, e.Status)
}

type Response struct {
	Name   string
	Price  float64
	Volume float64
}

type Pair struct {
	Base  string
	Quote string
}

type Exchange struct {
	pairs []*Pair
}

type Interface interface {
	GetConfig() *Config
	GetResponse(base, quote string) (*Response, error)
	RefreshPairs() error
	GetPairs() []*Pair
}

func GetSupportedExchanges() []Interface {
	return exchanges
}

func (exc *Exchange) GetPairs() []*Pair {
	return exc.pairs
}

func (exc *Exchange) SetPairs(pairs []*Pair) {
	exc.pairs = pairs
}

func (exc *Exchange) HttpGet(config *Config, url string, excModel interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.BaseURL, url), nil)
	if err != nil {
		return &Error{
			Exchange: config.Name,
			Status:   "500 ERROR",
			Message:  fmt.Sprintf("error on forming request to %s", config.Name)}
	}
	resp, err := client.Do(req)
	if err != nil {
		return &Error{Exchange: config.Name, Status: "500 ERROR", Message: err.Error()}
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Error{Exchange: config.Name, Status: resp.Status, Message: err.Error()}
	}
	err = json.Unmarshal(bodyBytes, excModel)
	if err != nil {
		return &Error{Exchange: config.Name, Status: resp.Status, Message: err.Error()}
	}
	return nil
}

func (exc *Exchange) toFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}

	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		vStr := v.(string)
		vF, _ := strconv.ParseFloat(vStr, 64)
		return vF
	default:
		panic("to float64 error.")
	}
}
