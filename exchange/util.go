package exchange

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var (
	exchanges = []Exchange{
		&Binance{},
		&Bitfinex{},
		//&Bitstamp{}, (Extreme rate limit)
		&Bittrex{},
		&GDAX{},
		&Huobi{},
		&HitBtc{},
		&ZB{} }
)

type Config struct {
	Name string
	BaseUrl string
	Client interface{}
	Pairs []*Pair
}

type Error struct {
	Exchange string
	Status string
	Message string
}

type Response struct {
	Name string
	Price float64
	Volume float64
}

type Pair struct {
	Base string
	Quote string
}

type Exchange interface {
	GetConfig() *Config
	GetPairs() []*Pair
	GetResponse(base, quote string) (*Response, *Error)
	SetPairs()
}

func GetSupportedExchanges() []Exchange {
	return exchanges
}

func HttpGet(config *Config, url string, excModel interface{}) *Error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.BaseUrl, url), nil)
	if err != nil {
		return &Error{
			Exchange: config.Name,
			Status: "500 ERROR",
			Message: fmt.Sprintf("error on forming request to %s", config.Name)}
	}
	resp, err := client.Do(req)
	if err != nil {
		return &Error{Exchange: config.Name, Status: resp.Status, Message: err.Error()}
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

func ToFloat64(v interface{}) float64 {
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
