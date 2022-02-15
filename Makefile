ifeq ($(OS),Windows_NT)
	COMPILE_TIME := $(shell echo %date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%)
else
	COMPILE_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
endif

.PHONY: default all

default: all

gsx2json:
	go build -ldflags '-s -w -X "main.Version=1.0.0" -X "main.Build=$(COMPILE_TIME)"' \
		-o ./build/ ./cmd/gsx2json/...

gencert:
	go build -o ./build/ ./cmd/gencert/...

install:
	go install cmd/gsx2json/gsx2json.go

fmt:
	go fmt ./...

tidy:
	go mod tidy

image:
	docker build -t gsx2json .

all: fmt gencert gsx2json tidy