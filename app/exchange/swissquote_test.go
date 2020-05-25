package exchange

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go test -run TestSwissquote_SetPairs
func TestSwissquote_SetPairs(t *testing.T) {
	swissquote := Swissquote{}
	_ = swissquote.RefreshPairs()
	pairs := swissquote.GetPairs()

	assert.Contains(t, pairs, &Pair{"XAU", "USD"})
	assert.Contains(t, pairs, &Pair{"XAG", "EUR"})
}

func TestSwissquote_GetResponse(t *testing.T) {
	swissquote := Swissquote{}
	price, err := swissquote.GetResponse("XAU", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Swissquote isn't greater than 0")
	//assert.True(t, price.Volume > 0, "volume from Swissquote isn't greater than 0")
}

// curl -X POST -H 'Content-Type: application/json' -d '{ "jobRunId": "1234", "data": { "base": "XAU", "quote": "USD" }}' http://localhost:8080/
