CMDS := .build/bin/houndd .build/bin/hound

SRCS := $(shell find . -type f -name '*.go')
UI := $(shell find ui/assets -type f)

WEBPACK_ARGS := --mode production
ifdef DEBUG
	WEBPACK_ARGS := --mode development
endif

ALL: $(CMDS)

ui: ui/.build/ui

# the mtime is updated on a directory when its files change so it's better
# to rely on a single file to represent the presence of node_modules.
node_modules/build:
	npm install
	date -u >> $@

.build/bin/houndd: ui/.build/ui $(SRCS)
	go build -o $@ github.com/hound-search/hound/cmds/houndd

.build/bin/hound: $(SRCS)
	go build -o $@ github.com/hound-search/hound/cmds/hound

ui/.build/ui: node_modules/build $(UI)
	mkdir -p ui/.build/ui
	cp -r ui/assets/* ui/.build/ui
	npx webpack $(WEBPACK_ARGS)

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
	rm -rf .build ui/.build node_modules
