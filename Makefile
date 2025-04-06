.PHONY: build run debug

BUILD_DATE := $(shell date -u +%Y-%m-%d)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/MWT-proger/time-tracking/internal/app.BuildDate=$(BUILD_DATE) -X github.com/MWT-proger/time-tracking/internal/app.GitCommit=$(GIT_COMMIT)

build:
	go build -ldflags "$(LDFLAGS)" -o time-tracker ./cmd/time-tracker

run: build
	./time-tracker -data ./data.json -notify-time 5

debug: build
	./time-tracker -data ./data.json -notify-time 5 -log-level debug 