FROM golang:1.14 as build

ARG WORKDIR=/go/src/github.com/golang-oauth2-example/

WORKDIR "${WORKDIR}"

COPY go.mod go.sum "${WORKDIR}"
RUN go mod download

COPY . "${WORKDIR}"

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o /example-oauth2-app

ENTRYPOINT ["/example-oauth2-app"]

FROM scratch

WORKDIR /app

# Copy the binary we built in the 'build' stage
COPY --from=build /example-oauth2-app .

# Copy across the CA certificates so that we can make TLS/SSL connections to things
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/app/example-oauth2-app"]
