package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	InitialiseConfig()
}

func TestBTCUSD(t *testing.T) {
	response := getPairResponse("BTC", "USD")

	assert.True(t, response.Params.Price != "0", "price returned from API is 0")
	assert.True(t, response.Params.Volume != "0", "volume returned from API is 0")
	assert.True(t, len(response.Params.Exchanges) > 1, "exchanges returned from API is less than 2")
	assert.Equal(t, response.Params.Id, "BTC-USD", "id of trading pair isn't correct")

	assert.Equal(
		t,
		response.Params.Price,
		response.Params.USDPrice.String,
		"Price is meant to match USD price for USD quotes",
	)
}

func TestETHEUR(t *testing.T) {
	response := getPairResponse("ETH", "EUR")

	assert.True(t, response.Params.Price != "0", "price returned from API is 0")
	assert.True(t, response.Params.Volume != "0", "volume returned from API is 0")
	assert.True(t, len(response.Params.Exchanges) > 1, "exchanges returned from API is less than 2")
	assert.Equal(t, response.Params.USDPrice.String, "", "usd price returned from API is 0")
	assert.Equal(t, response.Params.Id, "ETH-EUR", "id of trading pair isn't correct")
}

func TestLINKETH(t *testing.T) {
	response := getPairResponse("LINK", "ETH")

	assert.True(t, response.Params.Price != "0", "price returned from API is 0")
	assert.True(t, response.Params.Volume != "0", "volume returned from API is 0")
	assert.True(t, response.Params.USDPrice.String != "0", "usd price is 0")
	assert.True(t, len(response.Params.Exchanges) > 0, "exchanges returned from API was 0")
	assert.Equal(t, response.Params.Id, "LINK-ETH", "id of trading pair isn't correct")
}

func TestREQBTC(t *testing.T) {
	response := getPairResponse("REQ", "BTC")

	assert.True(t, response.Params.Price != "0", "price returned from API is 0")
	assert.True(t, response.Params.Volume != "0", "volume returned from API is 0")
	assert.True(t, response.Params.USDPrice.String != "0", "usd price is 0")
	assert.True(t, len(response.Params.Exchanges) > 0, "exchanges returned from API was 0")
	assert.Equal(t, response.Params.Id, "REQ-BTC", "id of trading pair isn't correct")
}

func TestUnknownPair(t *testing.T) {
	response := getPairResponse("UNK", "UNK")
	assert.Equal(t, response.Params.Price, "", "price returned from API is 0")
	assert.Equal(t, response.Params.Volume, "", "volume returned from API is 0")
	assert.Equal(t, response.Params.USDPrice.String, "", "usd price returned from API is 0")
	assert.Equal(t, len(response.Params.Exchanges), 0, "exchanges returned from API was 0")
	assert.Equal(t, response.Params.Id, "", "id of trading pair isn't correct")
}

func getPairResponse(base, quote string) *RunResult {
	server := httptest.NewServer(Api().MakeHandler())
	defer server.Close()

	var runResult RunResult
	runResult.JobRunID = "1234"
	runResult.Params.Base = base
	runResult.Params.Quote = quote
	b, err := json.Marshal(&runResult)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(fmt.Sprintf("%s/price", server.URL), "application/json", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	priceBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := RunResult{}
	err = json.Unmarshal(priceBody, &response)
	if err != nil {
		log.Fatal(err)
	}

	return &response
}
