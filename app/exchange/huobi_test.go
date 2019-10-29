package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestHuobi_SetPairs(t *testing.T) {
	huobi := Huobi{}
	_ = huobi.RefreshPairs()
	pairs := huobi.GetPairs()

	assert.Contains(t, pairs, &Pair{"ETH", "USDT"})
	assert.Contains(t, pairs, &Pair{"ETH", "BTC"})
}

func TestHuobi_GetResponse(t *testing.T) {
	huobi := Huobi{}
	price, err := huobi.GetResponse("ETH", "USDT")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Huobi isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Huobi isn't greater than 0")
}
