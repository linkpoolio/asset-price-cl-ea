package exchange

import "net/http"

type SupportedExchanges struct {
	exchanges []Exchange
}

type Response struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type Config struct {
	BaseUrl string
	Client *http.Client
}

type Exchange interface {
	GetConfig() Config
	GetTicker(base, quote string) (*Response, error)
}

func GetSupportedExchanges() []Exchange {
	exchanges := []Exchange{ GDAX{} }
	return exchanges
}