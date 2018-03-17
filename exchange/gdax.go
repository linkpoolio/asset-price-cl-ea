package exchange

import (
	"github.com/preichenberger/go-gdax"
	"fmt"
)

type GDAX struct {
	Exchange
}

func (exchange GDAX) GetTicker(base, quote string) (*Response, error) {
	client := gdax.NewClient("", "", "")

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, err
	}

	return &Response{Price: ticker.Price, Volume: ticker.Volume}, nil
}