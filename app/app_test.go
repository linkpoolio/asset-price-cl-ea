package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	StartPairsTicker(&Config{TickerInterval: time.Minute})
}

func TestBTCUSD(t *testing.T) {
	r, err := GetPrice("BTC", "USD")
	assert.Nil(t, err)

	assert.True(t, r.Price != "0", "price returned from API is 0")
	assert.True(t, r.Volume != "0", "volume returned from API is 0")
	assert.True(t, len(r.Exchanges) > 1, "exchanges returned from API is less than 2")
	assert.Equal(t, r.ID, "BTC-USD", "id of trading pair isn't correct")

	assert.Equal(
		t,
		r.Price,
		r.USDPrice.String,
		"Price is meant to match USD price for USD quotes",
	)
}

func TestETHEUR(t *testing.T) {
	r, err := GetPrice("ETH", "EUR")
	assert.Nil(t, err)

	assert.True(t, r.Price != "0", "price returned from API is 0")
	assert.True(t, r.Volume != "0", "volume returned from API is 0")
	assert.True(t, len(r.Exchanges) > 1, "exchanges returned from API is less than 2")
	assert.True(t,  r.USDPrice.String != "0", "usd price returned from API is 0")
	assert.Equal(t, r.ID, "ETH-EUR", "id of trading pair isn't correct")
}

func TestLINKETH(t *testing.T) {
	r, err := GetPrice("LINK", "ETH")
	assert.Nil(t, err)

	assert.True(t, r.Price != "0", "price returned from API is 0")
	assert.True(t, r.Volume != "0", "volume returned from API is 0")
	assert.True(t, r.USDPrice.String != "0", "usd price is 0")
	assert.True(t, len(r.Exchanges) > 0, "exchanges returned from API was 0")
	assert.Equal(t, r.ID, "LINK-ETH", "id of trading pair isn't correct")
}

func TestREQBTC(t *testing.T) {
	r, err := GetPrice("REQ", "BTC")
	assert.Nil(t, err)

	assert.True(t, r.Price != "0", "price returned from API is 0")
	assert.True(t, r.Volume != "0", "volume returned from API is 0")
	assert.True(t, r.USDPrice.String != "0", "usd price is 0")
	assert.True(t, len(r.Exchanges) > 0, "exchanges returned from API was 0")
	assert.Equal(t, r.ID, "REQ-BTC", "id of trading pair isn't correct")
}

func TestUnknownPair(t *testing.T) {
	_, err := GetPrice("UNK", "UNK")
	assert.Equal(t, fmt.Sprintf("%v", err), "No exchanges support that trading pair")
}
