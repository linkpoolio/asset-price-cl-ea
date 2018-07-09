package exchange

import (
	"github.com/adshao/go-binance"
	"context"
	"fmt"
	"log"
)

type Binance struct {
	Exchange
}

type BinanceProduct struct {
	LastPrice string `json:"lastprice"`
	Volume string `json:"quoteVolume"`
}

var binancePairs []*Pair

func (exchange *Binance) GetResponse(base, quote string) (*Response, *Error) {
	var bp BinanceProduct
	config := exchange.GetConfig()
	excErr := HttpGet(config, fmt.Sprintf("/ticker/24hr?symbol=%s%s", base, quote), &bp)
	if excErr != nil {
		return nil, excErr
	}
	return &Response{config.Name, ToFloat64(bp.LastPrice), ToFloat64(bp.Volume)}, nil
}

func (exchange *Binance) SetPairs() {
	clientInterface := exchange.GetConfig().Client
	client := clientInterface.(*binance.Client)

	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range exchangeInfo.Symbols {
		binancePairs = append(binancePairs, &Pair{product.BaseAsset, product.QuoteAsset})
	}
}

func (exchange *Binance) GetConfig() *Config {
	return &Config{
		Name: "Binance",
		BaseUrl: "https://www.binance.com/api/v1",
		Client: binance.NewClient("", ""),
		Pairs: binancePairs}
}
