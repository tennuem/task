FROM golang:1.12-alpine AS build

ENV PATH_ROJECT=${GOPATH}/src/github.com/tennuem/task
ENV APP=./cmd/app
ENV BIN=${GOPATH}/bin/task
ENV GO111MODULE=on
ARG VERSION
ENV VERSION ${VERSION:-0.1.0}
ARG BUILD_TIME
ENV BUILD_TIME ${BUILD_TIME:-unknown}
ARG COMMIT
ENV COMMIT ${COMMIT:-unknown}

RUN apk add --no-cache git
WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BIN} ${APP}

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
        -X github.com/tennuem/task/pkg/health.Version=${VERSION} \
        -X github.com/tennuem/task/pkg/health.Commit=${COMMIT} \
        -X github.com/tennuem/task/pkg/health.BuildTime=${BUILD_TIME}" \
        -a -installsuffix cgo -o ${BIN} ${APP}

FROM scratch
COPY --from=build /go/bin/task /bin/task
ENTRYPOINT ["/bin/task"]
