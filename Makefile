# enable buildkit
export DOCKER_BUILDKIT=1

# Some commands can be duplicated locally for better local development
# Set ENV to local to force DUPLICATE COMMANDS for local environment
ENV?=
ifeq (${ENV},)
	# Automatically ENV set up to local only if there is GO compiler
	ifeq ($(shell which go > /dev/null && echo "1" || echo "0"), 1)
	ENV=local
	endif
endif

# Builder cache ttl in hours
CACHE_TTL=1

all: build test check

.PHONY: build
build: modules build-dry

.PHONY: build-dry
build-dry:
	@docker build . \
 	--target scratch-out \
	--output ./bin

.PHONY: build-debug
build-debug:
	@docker build . \
	--build-arg DEBUG=1 \
 	--target scratch-out \
	--output ./bin


.PHONY: modules
modules:
	@docker build . \
 	--target go-mod-tidy \
	--output .

ifeq (${ENV}, local)
	go mod download
endif

.PHONY: test
test:
	@docker build . \
	--target test

.PHONY: check
check:
	@docker build . --target lint

.PHONY: cache-clear
cache-clear: ## Clear the builder cache
	docker builder prune --force --filter type=exec.cachemount --filter=unused-for==${CACHE_TTL}h
