ifeq ($(OS),Windows_NT)
	COMPILE_TIME := $(shell echo %date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%)
else
	COMPILE_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
endif

VERSION=1.0.3
BUILD=$(shell git rev-parse HEAD)
RELEASE=$(COMPILE_TIME)

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags '-s -w -X "main.Version=$(VERSION)" -X "main.Build=$(BUILD)" -X "main.Release=$(RELEASE)"'

# Sperate "linux-amd64" as GOOS and GOARCH
OSARCH_SPERATOR = $(word $2,$(subst -, ,$1))

.PHONY: default all

# Local build options
gsx2json: fmt tidy
	go build $(LDFLAGS) -o ./build/ ./cmd/gsx2json/...

gencert: fmt tidy
	go build $(LDFLAGS) -o ./build/ ./cmd/gencert/...

build: gsx2json gencert

# Architecture build options
gsx2json-arch-%: export GOARCH=$(call OSARCH_SPERATOR,$*,1)
gsx2json-arch-%: fmt tidy
	go build $(LDFLAGS) -o ./build/$(GOARCH)/ ./cmd/gsx2json/...

gencert-arch-%: export GOARCH=$(call OSARCH_SPERATOR,$*,1)
gencert-arch-%: fmt tidy
	go build $(LDFLAGS) -o ./build/$(GOARCH)/ ./cmd/gencert/...

arch-%: export GOARCH=$(call OSARCH_SPERATOR,$*,1)
arch-%: fmt tidy
	go build $(LDFLAGS) -o ./build/$(GOARCH)/ ./cmd/...

# Platform build options
cross-compile-%: export GOOS=$(call OSARCH_SPERATOR,$*,1)
cross-compile-%: export GOARCH=$(call OSARCH_SPERATOR,$*,2)
cross-compile-%: fmt tidy
	go build $(LDFLAGS) -o ./build/$(GOOS)-$(GOARCH)/ ./cmd/...

linux: cross-compile-linux-amd64
darwin: cross-compile-darwin-amd64
windows: cross-compile-windows-amd64

all: darwin linux windows

# Docker options
image:
	docker build -t gsx2json .

# Misc
fmt:
	go fmt ./...

tidy:
	go mod tidy

default: all