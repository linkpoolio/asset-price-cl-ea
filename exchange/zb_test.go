package exchange

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestZB_SetPairs(t *testing.T) {
	zb := ZB{}
	zb.SetPairs()
	pairs := zb.GetConfig().Pairs

	assert.Contains(t, pairs, &Pair{"ETH", "USDT"})
	assert.Contains(t, pairs, &Pair{"ETH", "BTC"})
}

func TestZB_GetResponse(t *testing.T) {
	zb := ZB{}
	price, err := zb.GetResponse("ETH", "USDT")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from ZB isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from ZB isn't greater than 0")
}