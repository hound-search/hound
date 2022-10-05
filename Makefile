CMDS := .build/bin/houndd .build/bin/hound

SRCS := $(shell find . -type f -name '*.go')

WEBPACK_ARGS := --mode production
ifdef DEBUG
	WEBPACK_ARGS := --mode development
endif

ALL: $(CMDS)

ui: ui/bindata.go

# the mtime is updated on a directory when its files change so it's better
# to rely on a single file to represent the presence of node_modules.
node_modules/build:
	npm install
	date -u >> $@

.build/bin/houndd: ui/bindata.go $(SRCS)
	go build -o $@ github.com/hound-search/hound/cmds/houndd

.build/bin/hound: ui/bindata.go $(SRCS)
	go build -o $@ github.com/hound-search/hound/cmds/hound

.build/bin/go-bindata:
	go build -o $@ github.com/go-bindata/go-bindata/go-bindata

ui/bindata.go: .build/bin/go-bindata node_modules/build $(wildcard ui/assets/**/*)
	rsync -r ui/assets/* .build/ui
	npx webpack $(WEBPACK_ARGS)
	$< -o $@ -pkg ui -prefix .build/ui -nomemcopy .build/ui/...

dev: node_modules/build

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
