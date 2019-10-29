package exchange

import (
	"fmt"
)

type Binance struct {
	Exchange
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

func (exc *Binance) GetResponse(base, quote string) (*Response, error) {
	var bp BinanceProduct
	config := exc.GetConfig()
	excErr := exc.HttpGet(config, fmt.Sprintf("/ticker/24hr?symbol=%s%s", base, quote), &bp)
	if excErr != nil {
		return nil, excErr
	}
	return &Response{config.Name, exc.toFloat64(bp.LastPrice), exc.toFloat64(bp.Volume)}, nil
}

func (exc *Binance) RefreshPairs() error {
	var bi BinanceInfo
	config := exc.GetConfig()

	err := exc.HttpGet(config, "/exchangeInfo", &bi)
	if err != nil {
		return err
	}

	var pairs []*Pair
	for _, product := range bi.Symbols {
		pairs = append(pairs, &Pair{product.BaseAsset, product.QuoteAsset})
	}
	exc.SetPairs(pairs)

	return nil
}

func (exc *Binance) GetConfig() *Config {
	return &Config{
		Name:    "Binance",
		BaseURL: "https://www.binance.com/api/v1",
	}
}
