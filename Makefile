VERSION := v0.0.4

ifndef GOBIN
ifndef GOPATH
$(error GOPATH is not set, please make sure you set your GOPATH correctly!)
endif
GOBIN=$(GOPATH)/bin
ifndef GOBIN
$(error GOBIN is not set, please make sure you set your GOBIN correctly!)
endif
endif

.PHONY: gen
gen: $(GOBIN)/mockery 
	@echo generating mocks...
	@go generate ./...

.PHONY: pre-push
pre-push: gen lint test

.PHONY: release
release: check-worktree lint test
	@VERSION=$(VERSION) ./hack/release.sh

.PHONY: check-worktree
check-worktree:
	@./hack/check-worktree.sh

PHONY: lint
lint: $(GOBIN)/golangci-lint
	@echo linting go code...
	@GOGC=off golangci-lint run --fix --timeout 6m
	
.PHONY: test
test:
	./hack/test.sh

.PHONY: clean
clean: 
	@rm -rf $(OUT_DIR) $(UI_OUT) vendor || true

./vendor:
	@echo vendoring...
	@go mod vendor

$(GOBIN)/mockery:
	@mkdir dist || true
	@echo installing: mockery
	@curl -L -o dist/mockery.tar.gz -- https://github.com/vektra/mockery/releases/download/v1.1.1/mockery_1.1.1_$(shell uname -s)_$(shell uname -m).tar.gz
	@tar zxvf dist/mockery.tar.gz mockery
	@chmod +x mockery
	@mkdir -p $(GOBIN)
	@mv mockery $(GOBIN)/mockery
	@mockery -version

$(GOBIN)/golangci-lint:
	@mkdir dist || true
	@echo installing: golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.49.0
