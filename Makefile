VERSION?=0.0.12
YEAR?=$(shell date +"%Y")

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: cur-version
cur-version:
	@echo -n $(VERSION)

.PHONY: lint
lint: $(GOBIN)/golangci-lint
	@golangci-lint run --fix --timeout 6m

.PHONY: codegen
codegen: $(GOBIN)/mockery
	go generate ./...
	go run ./hack/license.go --license ./hack/boilerplate.txt --year $(YEAR) .

.PHONY: test
test:
	@./hack/test.sh

.PHONY: tidy
tidy:
	@echo running go mod tidy...
	@go mod tidy

.PHONY: check-worktree
check-worktree:
	@./hack/check_worktree.sh

# noop - for ci
.PHONY: clean
clean:
	@echo cleaned 

.PHONY: release
release: tidy check-worktree fetch-tags
	./hack/release.sh

.PHONY: fetch-tags
fetch-tags:
	git fetch --tags

.PHONY: pre-push
pre-push:
	@make lint
	@make codegen
	@make check-worktree
	

$(GOBIN)/mockery:
	GO111MODULE=on go get github.com/vektra/mockery/v2@v2.5.1
	mockery --version

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.36.0
