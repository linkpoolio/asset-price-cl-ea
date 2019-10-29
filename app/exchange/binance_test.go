package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBinance_SetPairs(t *testing.T) {
	binance := Binance{}
	_ = binance.RefreshPairs()
	pairs := binance.GetPairs()

	assert.Contains(t, pairs, &Pair{"LINK", "ETH"})
	assert.Contains(t, pairs, &Pair{"REQ", "BTC"})
}

func TestBinance_GetResponse(t *testing.T) {
	binance := Binance{}
	price, err := binance.GetResponse("LINK", "ETH")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Binance isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Binance isn't greater than 0")
}
