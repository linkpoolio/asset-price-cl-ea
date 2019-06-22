package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCoinall_SetPairs(t *testing.T) {
	coinall := Coinall{}
	coinall.SetPairs()
	pairs := coinall.GetConfig().Pairs

	assert.Contains(t, pairs, &Pair{"ETH", "USDT"})
	assert.Contains(t, pairs, &Pair{"ETH", "BTC"})
}

func TestCoinall_GetResponse(t *testing.T) {
	coinall := Coinall{}
	price, err := coinall.GetResponse("ETH", "USDT")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Coinall isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Coinall isn't greater than 0")
}
