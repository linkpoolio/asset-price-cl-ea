# Asset Price External Adaptor ![Build](https://github.com/linkpoolio/asset-price-cl-ea/workflows/Build/badge.svg) [![codecov](https://codecov.io/gh/linkpoolio/asset-price-cl-ea/branch/master/graph/badge.svg)](https://codecov.io/gh/linkpoolio/asset-price-cl-ea)
External Adaptor for Chainlink which aggregates prices of crypto assets from multiple exchanges based on a weighted average of their volume.

This adaptor is built using the bridges framework: https://github.com/linkpoolio/bridges

### Currently Supported Exchanges:

- Binance
- Bitfinex
- Bitstamp (highly rate limited)
- Bittrex
- Coinall
- Coinbase Pro
- HitBTC
- Huobi Pro
- Kraken
- Gemini
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

#### AWS Lambda
1. Build the executable:
    ```bash
    GO111MODULE=on go build -o asset-price
    ```
2. Add the file to a ZIP archive:
    ```bash
    zip asset-price.zip ./asset-price
    ```
3. Upload the the zip file into AWS and then use `asset-price` as the
handler.
4. Set the `LAMBDA` environment variable to `true` in AWS for
the adaptor to be compatible with Lambda.

#### GCP Functions
1. Change into the app directory:
    ```bash
    cd app
    ```
2. Deploy into GCP
    ```bash
    gcloud functions deploy asset-price --runtime go111 --entry-point Handler --trigger-http
    ```

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
        "result": 3836.4042305857843,
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
        "result": 0.0031786459052877327,
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

### Solidity Usage

To use this adaptor on-chain, you can create the following Chainlink request:
```
Chainlink.Request memory req = buildChainlinkRequest(jobId, this, this.fulfill.selector);
req.add("base", "LINK");
req.add("quote", "BTC");
req.add("copyPath", "price");
req.addInt("times", 100000000);
```

Allowing you to change `base` and `quote` to any trading pair.

### Contribution
We welcome any contributors. The more exchanges supported, the better. Feel free to raise any PR's or issues.
