package exchange

type Config struct {
	Name string
	BaseUrl string
	Client interface{}
}

type Error struct {
	Exchange string
	Status string
	Message string
}

type Response struct {
	Name string
	Price float64
	Volume float64
}

type Pair struct {
	Base string
	Quote string
}

type SupportedExchanges struct {
	exchanges []Exchange
}

type Exchange interface {
	GetConfig() *Config
	GetResponse(base, quote string) (*Response, *Error)
	GetPairs() []*Pair
}

func GetSupportedExchanges() []Exchange {
	exchanges := []Exchange{ GDAX{}, Bitstamp{}, Binance{} }
	return exchanges
}

