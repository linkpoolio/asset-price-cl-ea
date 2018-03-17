# Asset Price External Adaptor
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
go get
go build -o asset-price-cl-ea
```

Then run the adaptor:
```
./asset-price-cl-ea
```

