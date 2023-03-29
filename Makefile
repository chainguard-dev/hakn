# Copyright 2023 Chainguard, Inc.
# SPDX-License-Identifier: Apache-2.0

KOCACHE_PATH ?= /tmp/ko
define create_kocache_path
  mkdir -p $(KOCACHE_PATH)
endef

GIT_VERSION ?= $(shell git describe --tags --always --dirty)

##########
# ko build
##########

.PHONY: ko
ko: ## Build the images and push to a registry
	$(create_kocache_path)
	ko resolve --platform=linux/amd64,linux/arm64 \
	--image-refs=serving.images \
	--tags ${GIT_VERSION} \
	-BRf config/serving > serving.yaml

.PHONY: ko-local
ko-local: ## Build the images and store locally for development
	$(create_kocache_path)
	ko resolve --platform=linux/amd64,linux/arm64 \
	--local --image-refs=serving.images \
	--tags ${GIT_VERSION} \
	-BRf config/serving > serving-local.yaml

##########
# release
##########

.PHONY: build-sign-image
build-sign-image: ko ## Used by Goreleaser to build and sign the images
	./hack/sign-images.sh

.PHONY: goreleaser
goreleaser: ## Run goreleaser to build the release images and create Github Relase
	goreleaser release --clean

##################
# help
##################

.PHONY: help
help: # Display help
	@awk -F ':|##' \
		'/^[^\t].+?:.*?##/ {\
			printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
