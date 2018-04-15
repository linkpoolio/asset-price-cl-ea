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
make build
```

Then run the adaptor:
```
export PORT=8080
./asset-price-cl-ea
```

#### Docker
To run the container:
```
docker run -it -p 8080:8080 linkpoolio/asset-price-cl-ea
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
    "price": 8298.837991560862,
    "volume": 16594.58795541,
    "exchanges": [
      "GDAX",
      "Bitstamp"
    ]
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
    "price": 0.00079999,
    "volume": 440625,
    "exchanges": [
      "Binance"
    ]
  },
  "status": "",
  "error": null,
  "pending": false
}
```

### ChainLink Node Setup

To integrate this adaptor with your node, use the following commands:

**Add Bridge Type**
```
curl -u <username>:<password> -X POST -H 'Content-Type: application/json' -d '{"name":"asset-price","url":"http://localhost:8080/price"}' http://localhost:6688/v2/bridge_types
```

**Create Spec**
```
curl -u <username>:<password> -X POST -H 'Content-Type: application/json' -d '{"initiators":[{"type":"web"}],"tasks":[{"type":"asset-price"},{"type":"noop"}]}' http://localhost:6688/v2/specs
```

**New Spec Run**

Notice the parameters `base` and `quote`. These are passed into the external adaptor by the node.
```
curl -u <username>:<password> -X POST -H 'Content-Type: application/json' -d '{"base": "BTC", "quote": "USD"}' http://localhost:6688
/v2/specs/<specId>/runs
```

### Contribution
We welcome any contributors. The more exchanges supported, the better. Feel free to raise any PR's or issues.