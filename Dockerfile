FROM golang:1.14 as build

ARG WORKDIR=/go/src/github.com/golang-oauth2-example/

WORKDIR "${WORKDIR}"

COPY go.mod go.sum "${WORKDIR}"
RUN go mod download

COPY . "${WORKDIR}"
RUN go build -o /example-oauth2-app

ENTRYPOINT ["/example-oauth2-app"]
