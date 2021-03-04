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

lint-%: $(GOBIN)/golangci-lint
	cd $(subst :,/,$*); golangci-lint run --fix --timeout 3m

vet-%:
	cd $(subst :,/,$*); go vet ./...

generate-%: $(GOBIN)/interfacer $(GOBIN)/mockery
	cd $(subst :,/,$*); go generate ./...; go run ../hack/license.go --license $(root_dir)/hack/boilerplate.txt --year $(YEAR) $(root_dir)

test-%:
	cd $(subst :,/,$*); $(root_dir)/hack/test.sh

pkg-%: generate-% tidy-% lint-% vet-% test-%;

coverage.txt:
	$(root_dir)/hack/test_pkg.sh

$(GOBIN)/mockery:
	GO111MODULE=on go get github.com/vektra/mockery/v2@v2.5.1
	mockery --version

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.36.0

$(GOBIN)/interfacer:
	GO111MODULE=on go get github.com/rjeczalik/interfaces/cmd/interfacer@v0.1.1