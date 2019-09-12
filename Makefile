DC = docker-compose
DK = docker
BIN = app
GO = go
CPROFILE = count.out
REPORT = report.xml
DKNAME = eikoapp/eiko
DKTAG = latest-prod

all: build-go-light
all: lint
all: test
all: clean

lint:
	golint

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
	$(GO) test -covermode=count -cover -coverprofile=$(CPROFILE) ./...

test-report:
	$(GO) test -covermode=count -cover -coverprofile=$(CPROFILE) ./... -v 2>&1 | go-junit-report > $(REPORT)

cover:
	$(GO) tool cover -html=$(CPROFILE) -o test.html

vet:
	$(GO) vet $(ARGS) ./...

clean:
	$(RM) $(BIN) $(CPROFILE) $(REPORT)

.PHONY: clean all build cover test