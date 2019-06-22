package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCoinbase_SetPairs(t *testing.T) {
	gdax := Coinbase{}
	gdax.SetPairs()
	pairs := gdax.GetConfig().Pairs

	assert.Contains(t, pairs, &Pair{"BTC", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "EUR"})
}

func TestCoinbase_GetResponse(t *testing.T) {
	gdax := Coinbase{}
	price, err := gdax.GetResponse("BTC", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Coinbase isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Coinbase isn't greater than 0")
}
