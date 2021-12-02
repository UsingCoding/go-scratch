export APP_CMD_NAME = scratch

APP_EXECUTABLE_OUT?=bin
BINARY?="${APP_EXECUTABLE_OUT}/${APP_CMD_NAME}"

# Empty string if there is no commit
GIT_COMMIT?=$(shell git rev-parse --short HEAD 2> /dev/null)

LDFLAGS="-X main.commit=${GIT_COMMIT}"

STATIC_FLAGS=CGO_ENABLED=0

GO_BUILD=$(STATIC_FLAGS) go build -trimpath -v -ldflags=$(LDFLAGS)

ifdef DEBUG
	GO_BUILD+="--gcflags=\"all=-N -l\""
	BINARY := "${BINARY}-debug"
endif

.PHONY: build
build:
	${GO_BUILD} -o ${BINARY} ./cmd/${APP_CMD_NAME}

.PHONY: modules
modules:
	go mod tidy

.PHONY: check
check:
	golangci-lint run

.PHONY: test
test:
	go test ./...
