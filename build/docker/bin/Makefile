SHELL = /bin/bash
VERSION ?= devel
GITCOMMIT = $(shell cd /src && git describe --always --dirty)
BUILDTIME = $(shell date --iso-8601=seconds)
LDFLAGS := -X almightyhelp/blockbook/blockbook/common.version=$(VERSION) -X almightyhelp/blockbook/blockbook/common.gitcommit=$(GITCOMMIT) -X almightyhelp/blockbook/blockbook/common.buildtime=$(BUILDTIME)
BLOCKBOOK_BASE := $(GOPATH)/src/almightyhelp/blockbook
BLOCKBOOK_SRC := $(BLOCKBOOK_BASE)/blockbook
ARGS ?=

all: build tools

build: prepare-sources generate-data
	cd $(BLOCKBOOK_SRC) && go build -o $(CURDIR)/blockbook -ldflags="-s -w $(LDFLAGS)" $(ARGS)
	cp $(CURDIR)/blockbook /out/blockbook
	chown $(PACKAGER) /out/blockbook

build-debug: prepare-sources generate-data
	cd $(BLOCKBOOK_SRC) && go build -o $(CURDIR)/blockbook -ldflags="$(LDFLAGS)" $(ARGS)
	cp $(CURDIR)/blockbook /out/blockbook
	chown $(PACKAGER) /out/blockbook

tools:
	cp $(CURDIR)/{ldb,sst_dump} /out
	chown $(PACKAGER) /out/{ldb,sst_dump}

test: prepare-sources generate-data
	cd $(BLOCKBOOK_SRC) && go test -tags unittest `go list ./... | grep -vP '^almightyhelp/blockbook/blockbook/(contrib|tests)'` $(ARGS)

test-integration: prepare-sources generate-data
	cd $(BLOCKBOOK_SRC) && go test -tags integration `go list almightyhelp/blockbook/blockbook/tests/...` $(ARGS)

test-all: prepare-sources generate-data
	cd $(BLOCKBOOK_SRC) && go test -tags 'unittest integration' `go list ./... | grep -v '^almightyhelp/blockbook/blockbook/contrib'` $(ARGS)

prepare-sources:
	@ [ -n "`ls /src 2> /dev/null`" ] || (echo "/src doesn't exist or is empty" 1>&2 && exit 1)
	rm -rf $(BLOCKBOOK_SRC)
	mkdir -p $(BLOCKBOOK_BASE)
	cp -r /src $(BLOCKBOOK_SRC)
	cd $(BLOCKBOOK_SRC) && go mod download

generate-data:
	cd $(BLOCKBOOK_SRC) && packr clean && packr
