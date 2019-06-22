package main

import (
	"github.com/linkpoolio/asset-price-cl-ea/app"
	"github.com/linkpoolio/bridges/bridge"
)

type AssetPrice struct{}

func (ap *AssetPrice) Run(h *bridge.Helper) (interface{}, error) {
	return app.GetPrice(h.GetParam("base"), h.GetParam("quote"))
}

func (ap *AssetPrice) Opts() *bridge.Opts {
	return &bridge.Opts{
		Name:   "Asset Price",
		Lambda: true,
		Path:   "/price",
	}
}

func main() {
	c := app.NewConfig()
	app.StartPairsTicker(c)
	bridge.NewServer(&AssetPrice{}).Start(c.Port)
}
