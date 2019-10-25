package exchange

import (
	"fmt"
)

type Binance struct {
	Exchange
	Pairs []*Pair
}

type BinanceProduct struct {
	LastPrice string `json:"lastprice"`
	Volume    string `json:"quoteVolume"`
}

type BinanceSymbol struct {
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}

type BinanceInfo struct {
	Symbols []*BinanceSymbol `json:"symbols"`
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

func (exc *Binance) SetPairs() *Error {
	var bi BinanceInfo
	config := exc.GetConfig()

	err := HttpGet(config, "/exchangeInfo", &bi)
	if err != nil {
		return err
	}

	for _, product := range bi.Symbols {
		exc.Pairs = append(exc.Pairs, &Pair{product.BaseAsset, product.QuoteAsset})
	}
	return nil
}

func (exc *Binance) GetConfig() *Config {
	return &Config{
		Name:    "Binance",
		BaseUrl: "https://www.binance.com/api/v1",
		Pairs:   exc.Pairs}
}
