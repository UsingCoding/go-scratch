# syntax=docker/dockerfile:1.2

ARG GO_VERSION=1.18
ARG GOLANGCI_LINT_VERSION=v1.43.0

FROM golang:${GO_VERSION} AS base
WORKDIR /app

RUN apt-get install \
    make

COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} AS lint-base

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    make -f rules/builder.mk check

FROM base AS test
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make -f rules/builder.mk test


FROM base as make-scracth

ARG DEBUG

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    APP_EXECUTABLE_OUT=/out \
    DEBUG=${DEBUG} \
    make -f rules/builder.mk

FROM scratch as scratch-out
COPY --from=make-scracth /out/* .

FROM base AS make-go-mod-tidy
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy

FROM scratch AS go-mod-tidy
COPY --from=make-go-mod-tidy /app/go.mod .
COPY --from=make-go-mod-tidy /app/go.sum .
