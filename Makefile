.DEFAULT_GOAL := build
.PHONY: build install docker dockerpush

REPO=linkpoolio/asset-price-cl-ea
LDFLAGS=-ldflags "-X github.com/linkpoolio/asset-price-cl-ea/store.Sha=`git rev-parse HEAD`"

gomod:
    export GO111MODULE=on

build: gomod
	@go build $(LDFLAGS) -o asset-price-cl-ea

install: gomod
	@go install $(LDFLAGS)

docker:
	@docker build . -t $(REPO)

dockerpush:
	@docker push $(REPO)