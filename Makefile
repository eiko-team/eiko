DC = docker-compose
DK = docker
BIN = app
GO = go
CPROFILE = count.out
COVFILE = coverage.txt
REPORT = report.xml
DKNAME = eikoapp/eiko
DKTAG = latest-prod

all: build-go-light
all: lint
all: test
all: clean

lint:
	golint ./...

build-go-light:
	$(GO) build -o $(BIN)

build-go:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GO) build -a -installsuffix cgo -o $(BIN)

build-dc:
	$(DC) build

build: build-go
build:
	$(DK) build -t "$(DKNAME):$(shell git rev-parse --short HEAD)" .

tag:
	$(DK) tag "$(DKNAME):$(shell git rev-parse --short HEAD)" "$(DKNAME):$(DKTAG)"

push:
	$(DK) push "$(DKNAME):$(shell git rev-parse --short HEAD)"

push-tag:
	$(DK) push "$(DKNAME):$(DKTAG)"

push-all: push
push-all: push-tag

up: build-go-light
up:
	$(DC) up

test:
	$(GO) test -tags mock -covermode=count -cover -coverprofile=$(CPROFILE) ./...

test-simple:
	$(GO) test -tags mock ./...

test-report:
	$(GO) test -tags mock -covermode=count -cover -coverprofile=$(CPROFILE) ./... -v 2>&1 | go-junit-report > $(REPORT)

test-full: test
test-full: test-simple

cover:
	$(GO) tool cover -html=$(CPROFILE) -o test.html

cover-race:
	$(GO) test -tags mock -race -coverprofile=$(COVFILE) -covermode=atomic ./...

vet:
	$(GO) vet $(ARGS) ./...

codecov: cover-race
codecov:
	$(shell curl -s https://codecov.io/bash | bash)

clean:
	$(RM) $(BIN) $(CPROFILE) $(REPORT) $(COVFILE)

.PHONY: clean all build cover test