APP?=./cmd/app
BIN?=./bin/auth

.PHONY: clean
clean:
	rm -f ${BIN}

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -a -installsuffix cgo -o ${BIN} ${APP}

.PHONY: run
run: build
	${BIN}

.PHONY: test
test:
	go test -v -race ./...