version := $(shell git describe --tags)
revision := $(shell git rev-parse HEAD)
release := $(shell git describe --tags | cut -d"-" -f 1,2)
IMPORTS := $(shell go list -f '{{join .Imports "\n"}}' ./... | sort | uniq | grep -v mtail)
GO_LDFLAGS := "-X main.Version=${version} -X main.Revision=${revision}"

.PHONY: install mailbeat
install mailbeat: $(GOFILES)
	go install -ldflags $(GO_LDFLAGS)

.PHONY: install_deps
install_deps:
	go get -u -v $(IMPORTS)