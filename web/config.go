package web

import (
	"time"
	"flag"
)

type Cfg struct {
	Port 		   int
	TickerInterval time.Duration
}

var Config Cfg

func InitialiseConfig() {
	Config = Cfg{}
	flag.IntVar(&Config.Port, "p", 8080, "Port number to serve")
	flag.DurationVar(
		&Config.TickerInterval,
		"t",
		time.Minute,
		"Ticker interval for the adaptor to refresh supported trading pairs, suggested units: s, m, h")
	flag.Parse()
}