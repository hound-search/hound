CMDS := $(GOPATH)/bin/houndd $(GOPATH)/bin/hound

SRCS := $(shell find . -type f -name '*.go')

WEBPACK_ARGS := -p
ifdef DEBUG
	WEBPACK_ARGS := -d
endif

ALL: $(CMDS)

ui: ui/.build/ui

node_modules:
	npm install

$(GOPATH)/bin/houndd: ui/.build/ui $(SRCS)
	go install github.com/hound-search/hound/cmds/houndd

$(GOPATH)/bin/hound: ui/bindata.go $(SRCS)
	go install github.com/hound-search/hound/cmds/hound

ui/.build/ui: node_modules $(wildcard ui/assets/**/*)
	mkdir -p ui/.build
	rsync -r ui/assets/* ui/.build/ui
	npx webpack $(WEBPACK_ARGS)

dev: ALL
	npm install

test:
	go test github.com/hound-search/hound/...
	npm test

lint:
	export GO111MODULE=on
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	export GOPATH=/tmp/gopath
	export PATH=$GOPATH/bin:$PATH
	golangci-lint run ./...

clean:
	rm -rf .build node_modules
