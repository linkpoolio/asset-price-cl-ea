package exchange

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestBitstamp_GetPairs(t *testing.T) {
	t.Parallel()

	bitstamp := Bitstamp{}
	pairs := bitstamp.GetPairs()

	assert.Contains(t, pairs, &Pair{"BTC", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "EUR"})
}

func TestBitstamp_GetPrice(t *testing.T) {
	t.Parallel()

	bitstamp := Bitstamp{}
	price, err := bitstamp.GetPrice("BTC", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Bitstamp isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Bitstamp isn't greater than 0")
}