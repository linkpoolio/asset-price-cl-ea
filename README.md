# Asset Price External Adaptor ![Travis-CI](https://travis-ci.org/linkpoolio/asset-price-cl-ea.svg?branch=master) [![codecov](https://codecov.io/gh/linkpoolio/asset-price-cl-ea/branch/master/graph/badge.svg)](https://codecov.io/gh/linkpoolio/asset-price-cl-ea)
External Adaptor for Chainlink which aggregates prices of crypto assets from multiple exchanges based on a weighted average of their volume.

### Currently Supported Exchanges:

- Binance
- Bitfinex
- ~~Bitstamp~~ (Disabled due to rate limit)
- Bittrex
- Coinall
- Coinbase Pro
- COSS
- HitBTC
- Huobi Pro
- ZB

### Setup Instructions
#### Local Install
Make sure [Golang](https://golang.org/pkg/) is installed.

Build:
```
make build
```

Then run the adaptor:
```
./asset-price-cl-ea -p <port> -t <tickerInterval>
```

##### Arguments

| Char   | Default  | Usage |
| ------ |:--------:| ----- |
| p      | 8080     | Port number to serve |
| t      | 1m0s     | Ticker interval for the adaptor to refresh supported trading pairs, suggested units: s, m, h |

#### Docker
To run the container:
```
docker run -it -p 8080:8080 -e PORT=8080 linkpoolio/asset-price-cl-ea
```

Container also supports passing in CLI arguments.

### Usage

To call the API, you need to send a POST request to `http://localhost:<port>/price` with the request body being of the ChainLink `RunResult` type.

For example:
```
curl -X POST -H 'Content-Type: application/json' -d '{ "jobRunId": "1234", "data": { "base": "BTC", "quote": "USD" }}' http://localhost:8080/price
```
Should return something similar to:
```json
{
    "jobRunId": "1234",
    "data": {
        "base": "BTC",
        "quote": "USD",
        "id": "BTC-USD",
        "price": "3836.4042305857843",
        "volume": "131747894.87525243",
        "usdPrice": "3836.4042305857843",
        "exchanges": [
            "HitBTC",
            "Bitfinex",
            "Coinbase",
            "COSS"
        ],
        "warnings": null
    },
    "status": "",
    "error": null,
    "pending": false
}
```

Or:
```
curl -X POST -H 'Content-Type: application/json' -d '{ "jobRunId": "1234", "data": { "base": "LINK", "quote": "ETH" }}' http://localhost:8080/price
```
```json
{
    "jobRunId": "1234",
    "data": {
        "base": "LINK",
        "quote": "ETH",
        "id": "LINK-ETH",
        "price": "0.0031786459052877327",
        "volume": "797.6642187877999",
        "usdPrice": "0.43956635389465454",
        "exchanges": [
            "Binance",
            "Huobi",
            "COSS"
        ],
        "warnings": null
    },
    "status": "",
    "error": null,
    "pending": false
}
```

### ChainLink Node Setup

To integrate this adaptor with your node, follow the official documentation:
https://docs.chain.link/docs/node-operators

### Contribution
We welcome any contributors. The more exchanges supported, the better. Feel free to raise any PR's or issues.
