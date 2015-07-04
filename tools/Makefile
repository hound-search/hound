ALL: commands

ASSETS=$(shell find src/github.com/etsy/Hound/ui/assets)

bin/go-bindata:
	GOPATH=$(shell pwd) go get github.com/jteeuwen/go-bindata/...

src/github.com/etsy/Hound/ui/bindata.go: bin/go-bindata $(ASSETS)
	mkdir -p build/ui
	rsync -r --exclude '*.js' src/github.com/etsy/Hound/ui/assets/* build/ui
	jsx --no-cache-dir src/github.com/etsy/Hound/ui/assets/js build/ui/js
	$< -o $@ -pkg ui -prefix build/ui -nomemcopy build/ui/...

commands: src/github.com/etsy/Hound/ui/bindata.go
	GOPATH=$(shell pwd) go install github.com/etsy/hound/cmds/...

clean:
	rm -rf bin/hound bin/houndd build

test:
	GOPATH=$(shell pwd) go test github.com/etsy/hound/...
