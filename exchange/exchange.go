package exchange

type Config struct {
	Name string
	BaseUrl string
}

type Response struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type SupportedExchanges struct {
	exchanges []Exchange
}

type Exchange interface {
	GetConfig() *Config
	GetResponse(base, quote string) (*Response, error)
}

func GetSupportedExchanges() []Exchange {
	exchanges := []Exchange{ GDAX{}, Bitstamp{} }
	return exchanges
}

