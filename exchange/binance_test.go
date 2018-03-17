package exchange

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestBinance_SetPairs(t *testing.T) {
	t.Parallel()

	binance := Binance{}
	binance.SetPairs()
	pairs := binance.GetConfig().Pairs

	assert.Contains(t, pairs, &Pair{"LINK", "ETH"})
	assert.Contains(t, pairs, &Pair{"REQ", "BTC"})
}

func TestBinance_GetPrice(t *testing.T) {
	t.Parallel()

	binance := Binance{}
	price, err := binance.GetPrice("LINK", "ETH")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Bitstamp isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Bitstamp isn't greater than 0")
}