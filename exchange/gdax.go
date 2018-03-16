package exchange

import (
	"os"

	"github.com/preichenberger/go-gdax"
	"fmt"
)

type GDAX struct {
	Exchange
}

func (exchange GDAX) GetTicker(base, quote string) (*Response, error) {
	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client := gdax.NewClient(secret, key, passphrase)

	ticker, err := client.GetTicker(fmt.Sprintf("%s-%s", base, quote))
	if err != nil {
		return nil, err
	}

	return &Response{Price: ticker.Price, Volume: ticker.Volume}, nil
}