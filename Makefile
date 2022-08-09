######################################################
# vars
######################################################
GOLANGCI_VERSION = 1.40.0

######################################################
# misc
######################################################
out:
	@mkdir -p out

out/bin:
	@mkdir -p out/bin


######################################################
# setup
######################################################
.PHONY: download
download: ## downloads the dependencies
	go mod download -x

######################################################
# clean
######################################################
.PHONY: clean-bin
clean-bin: ## clean local binary folders
	@rm -rf bin testbin

.PHONY: clean-outputs
clean-outputs: ## clean output folders out, vendor
	@rm -rf out vendor api/proto/google api/proto/validate

.PHONY: clean
clean: clean-bin clean-outputs ## clean up everything

######################################################
# lint
######################################################
bin/golangci-lint: bin/golangci-lint-$(GOLANGCI_VERSION)
	@ln -sf golangci-lint-$(GOLANGCI_VERSION) $@

bin/golangci-lint-$(GOLANGCI_VERSION):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v$(GOLANGCI_VERSION)
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint out download ## lint all code with golangci-lint
	# TODO:
	#bin/golangci-lint run ./... --timeout 15m0s

######################################################
# test
######################################################
.PHONY: test
test: download ## run all tests
	go test -v -coverpkg=./... -coverprofile=coverage.cov ./...

######################################################
# test pipeline
######################################################
.PHONY: test_pipeline
test_pipeline: download ## run all tests inside pipeline
	./ci/scripts/test.sh

######################################################
# coverage
######################################################
.PHONY: coverage
coverage: download ## run the code coverage
	./ci/scripts/cov.sh

######################################################
# build
######################################################
.PHONY: build
build: download out/bin ## build all binaries
	CGO_ENABLED=0 go build -ldflags="-w -s" -o out/bin ./...

.PHONY: build-linux
build-linux: download out/bin ## build all binaries for linux
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o out/bin ./...

.PHONY: build-mac
build-linux: download out/bin ## build all binaries for linux
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -ldflags="-w -s" -o out/bin ./...

######################################################
# help
######################################################
.PHONY: help
help: ## show help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''