package main

import (
	"github.com/linkpoolio/asset-price-cl-ea/app"
	"github.com/linkpoolio/bridges"
)

func main() {
	c := app.NewConfig()
	app.StartPairsTicker(c)

	bridges.NewServer(&app.AssetPrice{}).Start(c.Port)
}
