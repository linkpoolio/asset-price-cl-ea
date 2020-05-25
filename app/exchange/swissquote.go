package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// var goldSilverData []map[string][]map[string]float64
// 		json.Unmarshal(data, &goldSilverData)
// 		bid := goldSilverData[0]["spreadProfilePrices"][0]["bid"]
// 		ask := goldSilverData[0]["spreadProfilePrices"][0]["ask"]

type Swissquote struct {
	Exchange
}

type SwissquoteProduct struct {
	LastPrice string `json:"lastprice"`
	Volume    string `json:"quoteVolume"`
}

type SwissquoteSymbol struct {
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}

type SwissquoteInfo struct {
	Symbols []*SwissquoteSymbol `json:"symbols"`
}

var netClient = &http.Client{
	Timeout: time.Second * 5,
}

func (exc *Swissquote) GetResponse(base, quote string) (*Response, error) {
	var bid float64
	var ask float64
	config := exc.GetConfig()
	url := config.BaseURL + fmt.Sprintf("/%s/%s", base, quote)
	response, err := netClient.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	} else {

		data, _ := ioutil.ReadAll(response.Body)
		var goldSilverData []map[string][]map[string]float64
		json.Unmarshal(data, &goldSilverData)
		bid = goldSilverData[0]["spreadProfilePrices"][0]["bid"]
		ask = goldSilverData[0]["spreadProfilePrices"][0]["ask"]
	}
	median := (bid + ask) / 2
	return &Response{config.Name, exc.toFloat64(median), exc.toFloat64("1")}, nil
}

func (exc *Swissquote) RefreshPairs() error {
	//Hard coded pairs
	pairs := []*Pair{&Pair{"XAU", "USD"}, &Pair{"XAG", "USD"}, &Pair{"XAU", "EUR"}, &Pair{"XAG", "EUR"}, &Pair{"XAG", "GBP"}, &Pair{"XAU", "GBP"}, &Pair{"XAU", "AUD"}, &Pair{"XAG", "AUD"}, &Pair{"XAG", "CHF"}, &Pair{"XAU", "CHF"}}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Swissquote) GetConfig() *Config {
	return &Config{
		Name:    "Swissquote",
		BaseURL: "https://forex-data-feed.swissquote.com/public-quotes/bboquotes/instrument",
	}
}
