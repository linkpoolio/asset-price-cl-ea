package exchange

type Config struct {
	Name string
	BaseUrl string
	SupportedBases []string
	SupportedQuotes []string
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

type SupportedExchanges struct {
	exchanges []Exchange
}

type Exchange interface {
	GetConfig() *Config
	GetResponse(base, quote string) (*Response, *Error)
}

func GetSupportedExchanges() []Exchange {
	exchanges := []Exchange{ GDAX{}, Bitstamp{} }
	return exchanges
}

