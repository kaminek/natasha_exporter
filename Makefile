NAME := natasha_exporter
EXECUTABLE := $(NAME)
PACKAGES ?= $(shell go list ./... | grep -v /vendor/ | grep -v /_tools/)
SOURCES ?= $(shell find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./_tools/*")

.PHONY: all
all: dep build

.PHONY: clean
clean:
	go clean -i ./...
	rm -rf bin/ $(DIST)/

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: lint
lint:

.PHONY: dep
dep:
	dep ensure -update

.PHONY: install
install: $(SOURCES)
	go install -v  ./cmd/$(NAME)

.PHONY: build
build: dep bin/$(EXECUTABLE)

bin/$(EXECUTABLE): $(SOURCES)
	go build -i -v  -o $@ ./cmd/$(NAME)

.PHONY: docs
docs:
	hugo -s docs/
