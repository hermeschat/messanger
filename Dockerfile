FROM golang:1.12.2 AS builder

ARG PROJECT

COPY . /go/src/${PROJECT}

WORKDIR /go/src/${PROJECT}

RUN go get -v

RUN go build -v -o /go/bin/app

FROM ubuntu:18.04 AS app

ARG PROJECT

COPY --from=builder /go/bin/app /

RUN chmod +x /app

ENTRYPOINT /app