# Asset Price External Adaptor ![Travis-CI](https://travis-ci.org/linkpoolio/asset-price-cl-ea.svg?branch=master) [![codecov](https://codecov.io/gh/linkpoolio/asset-price-cl-ea/branch/master/graph/badge.svg)](https://codecov.io/gh/linkpoolio/asset-price-cl-ea)
External Adaptor for Chainlink which aggregates prices of crypto assets from multiple exchanges based on a weighted average of their volume.

### Currently Supported Exchanges:

- GDAX
- Bitstamp
- Binance

### Setup Instructions
#### Local Install
Make sure [Golang](https://golang.org/pkg/) is installed.

Build:
```
dep ensure -v
go build -o asset-price-cl-ea
```

Then run the adaptor:
```
export PORT=8080
./asset-price-cl-ea
```

### Usage

You can call the API with the following URL:
```
http://localhost:8080/:base/:quote/
```

For example:
```
curl http://localhost:8080/price/BTC/USD
```
Should return something similar to:
```json
{
  "id": "BTC-USD",
  "price": 7790,
  "volume": 25334.77576343,
  "exchanges": [
    "Bitstamp",
    "GDAX"
  ]
}
```

Or:
```
curl http://localhost:8080/price/LINK/ETH
```
```json
{
  "id": "LINK-ETH",
  "price": 0.00067673,
  "volume": 523251,
  "exchanges": [
    "Binance"
  ]
}
```

### Contribution
We welcome any contributors. The more exchanges supported, the better. Feel free to raise any PR's or issues.