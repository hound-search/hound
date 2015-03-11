BINS = bin/houndd bin/hound
HOST = localhost:6080

TESTS = ansi \
		hound/codesearch/index \
		hound/codesearch/regexp \
		hound/client \
		hound/index \
		hound/vcs \
		hound/config

ALL: $(BINS)

bin/houndd : $(wildcard src/hound/**/*)
	GOPATH=`pwd` go build -o $@ src/hound/cmds/houndd/main.go

bin/hound : $(wildcard src/hound/**/*)
	GOPATH=`pwd` go build -o $@ \
		-ldflags "-X main.defaultHost $(HOST)" \
		src/hound/cmds/hound/*.go

test:
	GOPATH=`pwd` go test $(TESTS)
	
clean:
	rm -f $(BINS)
