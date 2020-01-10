APP?=./cmd/app
BIN?=./bin/task
VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: clean
clean:
	rm -f ${BIN}

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -ldflags "-s -w \
        -X github.com/tennuem/task/pkg/health.Version=${VERSION} \
        -X github.com/tennuem/task/pkg/health.Commit=${COMMIT} \
        -X github.com/tennuem/task/pkg/health.BuildTime=${BUILD_TIME}" \
    	-a -installsuffix cgo -o ${BIN} ${APP}

.PHONY: run
run: build
	${BIN}

.PHONY: test
test:
	go test -v -race ./...