package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBitfinex_SetPairs(t *testing.T) {
	bitfinex := Bitfinex{}
	_ = bitfinex.RefreshPairs()
	pairs := bitfinex.GetPairs()

	assert.Contains(t, pairs, &Pair{"ETH", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "BTC"})
}

func TestBitfinex_GetResponse(t *testing.T) {
	bitfinex := Bitfinex{}
	price, err := bitfinex.GetResponse("ETH", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Bitfinex isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Bitfinex isn't greater than 0")
}
