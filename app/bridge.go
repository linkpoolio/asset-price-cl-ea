package app

import (
	"github.com/linkpoolio/bridges"
	"net/http"
)

type AssetPrice struct{}

func (ap *AssetPrice) Run(h *bridges.Helper) (interface{}, error) {
	return GetPrice(h.GetParam("base"), h.GetParam("quote"))
}

func (ap *AssetPrice) Opts() *bridges.Opts {
	return &bridges.Opts{
		Name:   "Asset Price",
		Lambda: true,
		Path:   "/price",
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	StartPairsTicker(nil)
	bridges.NewServer(&AssetPrice{}).Handler(w, r)
}