PACKAGE ?= infestor
GO_FILES := $(shell { git ls-files; } | grep ".go$$")
LICENSED_FILES := $(shell { git ls-files; } | grep ".go$$")

OUT_DIR := workdir
COVER_FILE := $(OUT_DIR)/cover.out
TEST_FLAGS ?= all

GO := go

clean:
	rm -rf $(OUT_DIR)
.PHONY: clean

lint:
	golangci-lint run ./... --timeout 5m
.PHONY: lint

test:
	$(GO) test ./... -tags $(TEST_FLAGS)
.PHONY: test

test-license: $(LICENSED_FILES)
	@grep -vlz "$$(tr '\n' . < LICENSE_HEADER)" $^ && exit 1 || exit 0
.PHONY: test-license

test-all: lint test-license
.PHONY: test-all

cover:
	@mkdir -p $(dir $(COVER_FILE))
	$(GO) test -tags $(TEST_FLAGS) -coverprofile=$(COVER_FILE) ./...
	$(GO) tool cover -func=$(COVER_FILE)
.PHONY: cover

bench:
	$(GO) test -tags $(TEST_FLAGS) -bench=. ./...
.PHONY: bench

add-license: $(LICENSED_FILES)
	for x in $^; do tmp=$$(cat LICENSE_HEADER; sed -n '/^package \|^\/\/ *+build /,$$p' $$x); echo "$$tmp" > $$x; done
.PHONY: add-license

VERSION_TAG_CURRENT := $(shell git tag --list 'v*' --points-at HEAD | sort --version-sort | tr \~ - | tail -1)
VERSION_TAG_LATEST := $(shell git tag --list 'v*' | tr - \~ | sort --version-sort | tr \~ - | tail -1)
ifeq ($(VERSION_TAG_CURRENT),$(VERSION_TAG_LATEST))
	VERSION := $(VERSION_TAG_CURRENT)
endif

VERSION_HASH := $(shell git rev-parse --short HEAD)
VERSION_DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
ifeq ($(VERSION),)
	VERSION := "dev-$(VERSION_HASH)-$(VERSION_DATE)"
endif

ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

LDFLAGS := -ldflags "-X github.com/makerdao/infestor.Version=$(VERSION)"
