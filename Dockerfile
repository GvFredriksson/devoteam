FROM golang:1.18.2-bullseye AS builder

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
RUN apt update
ENV BIN="/usr/local/bin"
ENV GOBIN="/go/bin"
WORKDIR /usr/src/app
COPY go.sum go.mod ./
ADD ./app/* ./app/
RUN go mod tidy
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

CMD go run app/main.go
