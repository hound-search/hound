BINS = bin/houndd bin/hound
HOST = localhost:6080

TESTS = ansi \
				code.google.com/p/codesearch/index \
				code.google.com/p/codesearch/regexp \
				hound/client \
				hound/index

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