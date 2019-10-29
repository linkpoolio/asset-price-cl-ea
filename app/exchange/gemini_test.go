package exchange

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGemini_SetPairs(t *testing.T) {
	gemini := Gemini{}
	_ = gemini.RefreshPairs()
	pairs := gemini.GetPairs()

	assert.Contains(t, pairs, &Pair{"ETH", "USD"})
	assert.Contains(t, pairs, &Pair{"BTC", "USD"})
}

func TestGemini_GetResponse(t *testing.T) {
	gemini := Gemini{}
	price, err := gemini.GetResponse("ETH", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Gemini isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Gemini isn't greater than 0")
}
