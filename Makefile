
ALL: ui/bindata.go

.build/bin/go-bindata:
	GOPATH=`pwd`/.build go get github.com/jteeuwen/go-bindata/...

ui/bindata.go: .build/bin/go-bindata $(wildcard ui/assets/**/*)
	rsync -r --exclude '*.js' ui/assets/* .build/ui
	jsx --no-cache-dir ui/assets/js .build/ui/js
	$< -o $@ -pkg ui -prefix .build/ui -nomemcopy .build/ui/...

clean:
	rm -rf .build


BUILD_ID ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DOCKER_IMAGE := hound-dev:$(BUILD_ID)

VOLUMES := \
	-v $(CURDIR):/go/src/github.com/etsy/hound \
	-v $(CURDIR)/dist/bin:/go/bin \
	-v $(CURDIR)/dist/pkg:/go/pkg

build:
	docker build -t $(DOCKER_IMAGE) -f Dockerfile.build .

dist:
	mkdir dist/

binary: dist build
	docker run --rm $(VOLUMES) $(DOCKER_IMAGE)

shell: dist build
	docker run --rm -ti $(VOLUMES) $(DOCKER_IMAGE) bash

build-image: binary
	docker build -t etsy/hound:$(BUILD_ID) .

test-unit: build
	docker run --rm -ti $(VOLUMES) $(DOCKER_IMAGE) go test -v ./...
