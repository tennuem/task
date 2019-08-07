FROM golang:1.12-alpine AS build
ENV PROJECT=task
ENV PATH_ROJECT=${GOPATH}/src/github.com/tennuem/${PROJECT}
ENV APP=${GOPATH}/bin/${PROJECT}
ENV GO111MODULE=on

RUN apk add --no-cache git
WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}
RUN go mod download
RUN cd cmd && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP} .

FROM scratch
COPY --from=build /go/bin/task /bin/task
ENTRYPOINT ["/bin/task"]
