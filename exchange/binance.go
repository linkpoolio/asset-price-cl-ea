package exchange

import (
	"github.com/adshao/go-binance"
	"context"
	"fmt"
	"log"
)

type Binance struct {
	Exchange
	Pairs []*Pair
}

type BinanceProduct struct {
	LastPrice string `json:"lastprice"`
	Volume string `json:"quoteVolume"`
}

func (exc *Binance) GetResponse(base, quote string) (*Response, *Error) {
	var bp BinanceProduct
	config := exc.GetConfig()
	excErr := HttpGet(config, fmt.Sprintf("/ticker/24hr?symbol=%s%s", base, quote), &bp)
	if excErr != nil {
		return nil, excErr
	}
	return &Response{config.Name, ToFloat64(bp.LastPrice), ToFloat64(bp.Volume)}, nil
}

func (exc *Binance) SetPairs() {
	clientInterface := exc.GetConfig().Client
	client := clientInterface.(*binance.Client)

	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range exchangeInfo.Symbols {
		exc.Pairs = append(exc.Pairs, &Pair{product.BaseAsset, product.QuoteAsset})
	}
}

func (exc *Binance) GetConfig() *Config {
	return &Config{
		Name: "Binance",
		BaseUrl: "https://www.binance.com/api/v1",
		Client: binance.NewClient("", ""),
		Pairs: exc.Pairs}
}
