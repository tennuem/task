FROM golang:1.12-alpine AS build

ENV PATH_ROJECT=${GOPATH}/src/github.com/tennuem/task
ENV APP=./cmd/app
ENV BIN=${GOPATH}/bin/task
ENV GO111MODULE=on

RUN apk add --no-cache git
WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BIN} ${APP}

FROM scratch
COPY --from=build /go/bin/task /bin/task
ENTRYPOINT ["/bin/task"]
