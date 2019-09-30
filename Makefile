DC = docker-compose
DK = docker
BIN = app
GO = go
CPROFILE = count.out
COVFILE = coverage.txt
REPORT = report.xml
DKNAME = eikoapp/eiko
DKTAG = latest-prod
UGLY-JS = uglifyjs
UGLY-CSS = uglifycss
TRASHDIR = static/min
HTML_MIN_ARGS = --collapse-whitespace \
				--remove-comments \
				--remove-optional-tags \
				--remove-redundant-attributes \
				--remove-script-type-attributes \
				--remove-tag-whitespace \
				--use-short-doctype \
				--minify-css true \
				--minify-js true

all: create-dir
all: mini
all: build-go-light
all: lint
all: vet
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
	$(DK) build --no-cache -t "$(DKNAME):$(shell git rev-parse --short HEAD)" .

tag:
	$(DK) tag "$(DKNAME):$(shell git rev-parse --short HEAD)" "$(DKNAME):$(DKTAG)"

push:
	$(DK) push "$(DKNAME):$(shell git rev-parse --short HEAD)"

push-tag:
	$(DK) push "$(DKNAME):$(DKTAG)"

docker: mini-css
docker: build
docker: push
docker: push-tag

up: create-dir
up: mini
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
	$(GO) vet $(ARGS) -tags mock ./...

clean:
	$(RM) -r $(BIN) $(CPROFILE) $(REPORT) $(COVFILE) $(TRASHDIR)

create-dir:
	mkdir -p $(TRASHDIR)

mini: mini-css
mini: mini-html
mini: mini-js
mini: mini-img

mini-css:
	mkdir -p static/min/css
	$(UGLY-CSS) $(ARGS) static/css/eiko.css --debug --output static/min/css/eiko.css

HTML = $(shell ls static/html | grep html)

mini-html:
	mkdir -p static/min/html
	for file in $(HTML); do \
		html-minifier $(ARGS) $(HTML_MIN_ARGS) static/html/$$file \
			 -o static/min/html/$$file; \
	done

mini-js:
	cp -r static/js static/min

mini-img:
	cp -r static/img static/min

.PHONY: clean all build cover test
