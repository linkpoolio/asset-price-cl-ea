package exchange

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type Bitstamp struct {
	Exchange
}

type BitstampModel struct {
	Last string `json:"last"`
	Volume string `json:"volume"`
}

func (exchange Bitstamp) GetResponse(base, quote string) (*Response, *Error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s/ticker/%s%s", exchange.GetConfig().BaseUrl, base, quote), nil)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, "500 ERROR", "error on forming request to Bitstamp"}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	bitstampModel := &BitstampModel{}
	err = json.Unmarshal(bodyBytes, bitstampModel)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	currentPrice, err := strconv.ParseFloat(bitstampModel.Last, 64)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	currentVolume, err := strconv.ParseFloat(bitstampModel.Volume, 64)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	return &Response{Price: currentPrice, Volume: currentVolume}, nil
}

func (exchange Bitstamp) GetConfig() *Config {
	return &Config{
		BaseUrl: "https://www.bitstamp.net/api/v2",
		Name: "Bitstamp",
		SupportedBases: []string{"BTC", "BCH", "LTC"},
		SupportedQuotes: []string{"USD", "EUR"}}
}