
ALL: ui/bindata.go commands

go-bindata:
	go get github.com/jteeuwen/go-bindata/...

ui/bindata.go: go-bindata $(wildcard ui/assets/**/*)
	rsync -r --exclude '*.js' ui/assets/* .build/ui
	jsx --no-cache-dir ui/assets/js .build/ui/js
	../../../../bin/go-bindata -o $@ -pkg ui -prefix .build/ui -nomemcopy .build/ui/...

commands: ui/bindata.go
	go install github.com/etsy/hound/cmds/...

clean:
	rm -rf .build

test:
	go test github.com/etsy/hound/...
