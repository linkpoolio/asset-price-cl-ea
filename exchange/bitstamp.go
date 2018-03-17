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

func (exchange Bitstamp) GetResponse(base, quote string) (*Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s/ticker/%s%s", exchange.GetConfig().BaseUrl, base, quote), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bitstampModel := &BitstampModel{}
	err = json.Unmarshal(bodyBytes, bitstampModel)
	if err != nil {
		return nil, err
	}
	currentPrice, err := strconv.ParseFloat(bitstampModel.Last, 64)
	if err != nil {
		return nil, err
	}
	currentVolume, err := strconv.ParseFloat(bitstampModel.Volume, 64)
	if err != nil {
		return nil, err
	}
	return &Response{Price: currentPrice, Volume: currentVolume}, nil
}

func (exchange Bitstamp) GetConfig() *Config {
	return &Config{BaseUrl: "https://www.bitstamp.net/api/v2", Name: "Bitstamp"}
}