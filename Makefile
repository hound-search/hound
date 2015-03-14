
ALL: ui/bindata.go

build/bin/go-bindata:
	GOPATH=`pwd`/build go get github.com/jteeuwen/go-bindata/...

ui/bindata.go: build/bin/go-bindata $(wildcard ui/assets/**/*)
	rsync -r --exclude '*.js' ui/assets/* build/ui
	jsx --no-cache-dir ui/assets/js build/ui/js
	$< -o $@ -pkg ui -prefix build/ui -nomemcopy build/ui/...

clean:
	rm -f build
