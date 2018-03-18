package exchange

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
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
	return &Response{exchange.GetConfig().Name, currentPrice, currentVolume}, nil
}

func (exchange Bitstamp) SetPairs() {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s/trading-pairs-info/", exchange.GetConfig().BaseUrl), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var bPairs []BitstampPair
	err = json.Unmarshal(bodyBytes, &bPairs)
	if err != nil {
		log.Fatal(err)
	}

	for _, pair := range bPairs {
		if pair.Trading == "Enabled" {
			currencies := strings.Split(pair.Name, "/")
			bitstampPairs = append(bitstampPairs, &Pair{currencies[0], currencies[1]})
		}
	}
}

func (exchange Bitstamp) GetConfig() *Config {
	return &Config{BaseUrl: "https://bitstamp.net/api/v2", Name: "Bitstamp", Pairs: bitstampPairs}
}