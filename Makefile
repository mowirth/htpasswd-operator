export GOBIN=$(PWD)/bin

.PHONY: build clean default
default: build
build:
	@mkdir -p bin/
	@go install -v ./...
clean:
	@rm -rf bin/
lint: bin/golangci-lint
	bin/golangci-lint run --fix --timeout=600s
bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@