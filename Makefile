VER?=0.0.1
YEAR?=$(shell date +"%Y")
MODULES=$(shell find . -mindepth 2 -maxdepth 4 -type f -name 'go.mod' | cut -c 3- | sed 's|/[^/]*$$||' | sort -u | tr / :)
targets=$(addprefix pkg-, $(MODULES))
root_dir=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: $(targets) coverage.txt

tidy-%:
	cd $(subst :,/,$*); go mod tidy

lint-%: $(GOPATH)/bin/golangci-lint
	cd $(subst :,/,$*); golangci-lint run --fix --timeout 3m

vet-%:
	cd $(subst :,/,$*); go vet ./...

generate-%: $(GOPATH)/bin/interfacer
	cd $(subst :,/,$*); go generate ./...; go run ../hack/license.go --license $(root_dir)/hack/boilerplate.txt --year $(YEAR) $(root_dir)

test-%:
	cd $(subst :,/,$*); $(root_dir)/hack/test.sh

pkg-%: generate-% tidy-% lint-% vet-% test-%;

coverage.txt:
	$(root_dir)/hack/test_pkg.sh

$(GOPATH)/bin/mockery:
	@curl -L -o dist/mockery.tar.gz -- https://github.com/vektra/mockery/releases/download/v1.1.1/mockery_1.1.1_$(shell uname -s)_$(shell uname -m).tar.gz
	@tar zxvf dist/mockery.tar.gz mockery
	@chmod +x mockery
	@mkdir -p $(GOPATH)/bin
	@mv mockery $(GOPATH)/bin/mockery
	@mockery -version

$(GOPATH)/bin/golangci-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin v1.36.0

$(GOPATH)/bin/interfacer:
	GO111MODULE=on go get github.com/rjeczalik/interfaces/cmd/interfacer@v0.1.1