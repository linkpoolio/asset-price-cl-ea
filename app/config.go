package app

import (
	"time"
	"flag"
)

type Config struct {
	Port 		   int
	TickerInterval time.Duration
}

func NewConfig() *Config {
	c := Config{}
	flag.IntVar(&c.Port, "p", 8080, "Port number to serve")
	flag.DurationVar(
		&c.TickerInterval,
		"t",
		time.Minute,
		"Ticker interval for the adaptor to refresh supported trading pairs, suggested units: s, m, h")
	flag.Parse()
	return &c
}