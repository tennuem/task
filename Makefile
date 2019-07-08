PROJECT?=task
APP?=bin/${PROJECT}

clean:
	rm -f ${APP}
build: clean
	cd ./cmd && \
		CGO_ENABLED=0 go build -o ../${APP} .
run: build
	./${APP}
test:
	go test -v -race ./...