package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestHitBtc_SetPairs(t *testing.T) {
	hitBtc := HitBtc{}
	hitBtc.SetPairs()
	pairs := hitBtc.GetConfig().Pairs

	assert.Contains(t, pairs, &Pair{"ETH", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "BTC"})
}

func TestHitBtc_GetResponse(t *testing.T) {
	hitBtc := HitBtc{}
	price, err := hitBtc.GetResponse("ETH", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from HitBTC isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from HitBTC isn't greater than 0")
}
