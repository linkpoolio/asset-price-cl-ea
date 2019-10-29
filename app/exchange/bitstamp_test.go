package exchange

/* API has changed, commenting entirely because of that & rate limiting
func TestBitstamp_SetPairs(t *testing.T) {
	bitstamp := Bitstamp{}
	_ = bitstamp.RefreshPairs()
	pairs := bitstamp.GetPairs()

	assert.Contains(t, pairs, &Pair{"BTC", "USD"})
	assert.Contains(t, pairs, &Pair{"ETH", "EUR"})
}

/* Disabled due to rate limit
func TestBitstamp_GetResponse(t *testing.T) {
	bitstamp := Bitstamp{}
	price, err := bitstamp.GetResponse("BTC", "USD")
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, price.Price > 0, "price from Bitstamp isn't greater than 0")
	assert.True(t, price.Volume > 0, "volume from Bitstamp isn't greater than 0")
}
*/