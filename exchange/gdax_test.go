package exchange

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestGDAX_SetPairs(t *testing.T) {
	t.Parallel()

	gdax := GDAX{}
	gdax.SetPairs()
	pairs := gdax.GetConfig().Pairs

	assert.Contains(t, pairs, &Pair{"BTC", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "EUR"})
}

func TestGDAX_GetResponse(t *testing.T) {
	t.Parallel()

	gdax := GDAX{}
	price, err := gdax.GetResponse("BTC", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from GDAX isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from GDAX isn't greater than 0")
}