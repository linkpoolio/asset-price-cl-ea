FROM golang:1.12-alpine as builder

RUN apk add --no-cache make curl git gcc musl-dev linux-headers

ADD . /go/src/github.com/linkpoolio/asset-price-cl-ea
RUN cd /go/src/github.com/linkpoolio/asset-price-cl-ea && make build

# Copy into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/linkpoolio/asset-price-cl-ea/asset-price-cl-ea /usr/local/bin/

EXPOSE 8080
CMD ["asset-price-cl-ea"]