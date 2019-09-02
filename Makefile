DC = docker-compose
DK = docker
BIN = app
GO = go
CPROFILE = count.out
REPORT = report.xml
DKNAME = eikoapp/eiko

all: build-go-light
all: test
all: clean

build-go-light:
	$(GO) build -o $(BIN)

build-go:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GO) build -a -installsuffix cgo -o $(BIN)

build-dc:
	$(DC) build

build: build-go
build:
	$(DK) build -t "$(DKNAME):$(shell git rev-parse --short HEAD)" .

up: build-go
up:
	$(DC) up

test:
	$(GO) test -covermode=count -cover -coverprofile=$(CPROFILE) ./src

test-report:
	$(GO) test -covermode=count -cover -coverprofile=$(CPROFILE) ./src -v 2>&1 | go-junit-report > $(REPORT)

cover:
	$(GO) tool cover -html=$(CPROFILE) -o test.html

clean:
	$(RM) $(BIN) $(CPROFILE) $(REPORT)

.PHONY: clean all build cover test