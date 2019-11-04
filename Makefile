DC = docker-compose
DK = docker
BIN = app
GO = go
GITHASH = $(shell git rev-parse --short HEAD)
CPROFILE = count.out
COVFILE = coverage.txt
REPORT = report.xml
DKNAME = eikoapp/eiko
DKTAG = latest-prod
DKTAGHASH = $(GITHASH)
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

all: build-go-light
all: lint
all: fmt
all: vet
all: test
all: clean

lint:
	golint ./...

fmt:
	./scripts/gofmt.sh

build-go-light:
	$(GO) build -o $(BIN)

build-go:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GO) build -ldflags="-w -s" -a -installsuffix cgo -o $(BIN)

build-pi:
	GOARM=5 CGO_ENABLED=0 GOARCH=arm GOOS=linux $(GO) build -ldflags="-w -s" -a -installsuffix cgo -o $(BIN)

build-dc:
	$(DC) build

build-docker:
	$(DK) build --no-cache -t "$(DKNAME):$(DKTAGHASH)" .

tag:
	$(DK) tag "$(DKNAME):$(DKTAGHASH)" "$(DKNAME):$(DKTAG)"

push:
	$(DK) push "$(DKNAME):$(DKTAGHASH)"

push-tag:
	$(DK) push "$(DKNAME):$(DKTAG)"

docker: mini
docker: build-go
docker: build-docker
docker: push
docker: tag
docker: push-tag

docker-pi:
	make DKTAGHASH=arm-$(DKTAGHASH) mini
	make DKTAGHASH=arm-$(DKTAGHASH) build-pi
	make DKTAGHASH=arm-$(DKTAGHASH) build-docker
	make DKTAGHASH=arm-$(DKTAGHASH) push
	make DKTAGHASH=arm-$(DKTAGHASH) DKTAG=arm-$(DKTAG) tag
	make DKTAGHASH=arm-$(DKTAGHASH) DKTAG=arm-$(DKTAG) push-tag

up: create-dir
up: build-go-light
up:
	$(DC) up

test:
	$(GO) test $(ARGS) -tags mock -covermode=count -cover -coverprofile=$(CPROFILE) ./...

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
mini: mini-json
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

mini-json:
	cp -r static/json static/min

mini-img:
	cp -r static/img static/min

p: test
p:
	git push -v
p: docker
p: docker-pi

.PHONY: clean all build cover test
