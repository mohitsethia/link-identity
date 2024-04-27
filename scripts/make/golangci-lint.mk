LOCAL_BIN:=$(CURDIR)/bin

# Check global GOLANGCI-LINT
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG:=1.54.2

# Check local bin version
ifneq ($(wildcard $(GOLANGCI_BIN)),)
GOLANGCI_BIN_VERSION:=$(shell $(GOLANGCI_BIN) --version)
ifneq ($(GOLANGCI_BIN_VERSION),)
GOLANGCI_BIN_VERSION_SHORT:=$(shell echo "$(GOLANGCI_BIN_VERSION)"|sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_BIN_VERSION_SHORT:=0
endif
ifneq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_BIN_VERSION_SHORT)))"
GOLANGCI_BIN:=
endif
endif

# Check global bin version
ifneq (, $(shell which golangci-lint))
GOLANGCI_VERSION:=$(shell golangci-lint --version 2> /dev/null )
ifneq ($(GOLANGCI_VERSION),)
GOLANGCI_VERSION_SHORT:=$(shell echo "$(GOLANGCI_VERSION)"|sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_VERSION_SHORT:=0
endif
ifeq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_VERSION_SHORT)))"
GOLANGCI_BIN:=$(shell which golangci-lint)
endif
endif

# install golangci-lint binary
.PHONY: install-golangci-lint
.SILENT: install-golangci-lint
install-golangci-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info #Downloading golangci-lint v$(GOLANGCI_TAG))
	tmp=$$(mktemp -d) && cd $$tmp && pwd && go mod init temp && go get -d github.com/golangci/golangci-lint@v$(GOLANGCI_TAG) && \
		go build -mod mod -ldflags "-X 'main.version=$(GOLANGCI_TAG)' -X 'main.commit=test' -X 'main.date=test'" -o $(LOCAL_BIN)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint && \
		rm -rf $$tmp
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

# run diff golangci-lint
.PHONY: .golangci-lint
.SILENT: .golangci-lint
.golangci-lint: install-golangci-lint
	$(info #Running golangci-lint...)
	$(GOLANGCI_BIN) run --new-from-rev=origin/master --config=.golangci.yaml ./...

# run full golangci-lint
.PHONY: .golangci-lint-full
.SILENT: .golangci-lint-full
.golangci-lint-full: install-golangci-lint
	$(info #Running golangci-lint-full...)
	$(GOLANGCI_BIN) run --config=.golangci.yaml ./...
