FROM golang:1.14 as build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN go build -o /example-oauth2-app

ENTRYPOINT ["/example-oauth2-app"]

# FROM scratch

# COPY --from=build /example-oauth2-app /usr/local/bin/example-oauth2-app

# ENTRYPOINT ["/usr/local/bin/example-oauth2-app"]
