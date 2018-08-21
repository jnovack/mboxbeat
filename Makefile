version := $(shell git describe --tags)
revision := $(shell git rev-parse HEAD)
release := $(shell git describe --tags | cut -d"-" -f 1,2)

.PHONY: install mailbeat
install mailbeat: $(GOFILES)
	$(eval GO_LDFLAGS := "-X main.Version=${version} -X main.Revision=${revision}")
	go install -ldflags $(GO_LDFLAGS)

.PHONY: install_deps
install_deps:
	$(eval IMPORTS := $(shell go list -f '{{join .Imports "\n"}}' ./... | sort | uniq | grep -v mboxbeat))
	go get -u -v $(IMPORTS)

.PHONY: build
build:
	docker build -t mboxbeat .

.PHONY: quietbuild
quietbuild:
	docker build -q -t mboxbeat .

.PHONY: run
run:
	@docker run -it --rm -v `pwd`/test/mbox:/mbox mboxbeat /mbox

.PHONY: test
test: quietbuild run
