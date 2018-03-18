package web

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestBTCUSD(t *testing.T) {
	response := getPairResponse("BTC", "USD")

	assert.True(t, response.Price > 0, "price returned from API is 0")
	assert.True(t, response.Volume > 0, "volume returned from API is 0")
	assert.True(t, len(response.Exchanges) > 1, "exchanges returned from API is less than 2")
	assert.Equal(t, response.Id, "BTC-USD", "id of trading pair isn't correct")
}

func TestETHEUR(t *testing.T) {
	response := getPairResponse("ETH", "EUR")

	assert.True(t, response.Price > 0, "price returned from API is 0")
	assert.True(t, response.Volume > 0, "volume returned from API is 0")
	assert.True(t, len(response.Exchanges) > 1, "exchanges returned from API is less than 2")
	assert.Equal(t, response.Id, "ETH-EUR", "id of trading pair isn't correct")
}

func TestLINKETH(t *testing.T) {
	response := getPairResponse("LINK", "ETH")

	assert.True(t, response.Price > 0, "price returned from API is 0")
	assert.True(t, response.Volume > 0, "volume returned from API is 0")
	assert.True(t, len(response.Exchanges) > 0, "exchanges returned from API was 0")
	assert.Equal(t, response.Id, "LINK-ETH", "id of trading pair isn't correct")
}

func TestREQBTC(t *testing.T) {
	response := getPairResponse("REQ", "BTC")

	assert.True(t, response.Price > 0, "price returned from API is 0")
	assert.True(t, response.Volume > 0, "volume returned from API is 0")
	assert.True(t, len(response.Exchanges) > 0, "exchanges returned from API was 0")
	assert.Equal(t, response.Id, "REQ-BTC", "id of trading pair isn't correct")
}

func TestUnknownPair(t *testing.T) {
	response := getPairResponse("UNK", "UNK")
	assert.Equal(t, response.Price, float64(0), "price returned from API is 0")
	assert.Equal(t, response.Volume, float64(0), "volume returned from API is 0")
	assert.Equal(t, len(response.Exchanges), 0, "exchanges returned from API was 0")
	assert.Equal(t, response.Id, "", "id of trading pair isn't correct")
}

func getPairResponse(base, quote string) *Response {
	server := httptest.NewServer(Api().MakeHandler())
	defer server.Close()

	res, err := http.Get(fmt.Sprintf("%s/price/%s/%s", server.URL, base, quote))
	if err != nil {
		log.Fatal(err)
	}

	priceBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{}
	err = json.Unmarshal(priceBody, &response)
	if err != nil {
		log.Fatal(err)
	}

	return &response
}