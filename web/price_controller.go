package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/linkpoolio/asset-price-cl-ea/exchange"
	"net/http"
	"fmt"
)

func GetPrice(w rest.ResponseWriter, r *rest.Request) {
	base := r.PathParam("base")
	quote := r.PathParam("quote")

	// Get responses from each exchange based on the inputs
	var responses []*exchange.Response
	for _, exc := range exchange.GetSupportedExchanges() {
		response, err := exc.GetTicker(base, quote)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responses = append(responses, response)
	}

	// Calculate the weighted average based on volume
	var totalVolume float64
	for _, response := range responses {
		totalVolume += response.Volume
	}
	var weightedPrice float64
	for _, response := range responses {
		weightedPrice += (response.Volume / totalVolume) * response.Price
	}
	w.WriteJson(&exchange.Response{fmt.Sprintf("%s-%s", base, quote), weightedPrice, totalVolume})
}