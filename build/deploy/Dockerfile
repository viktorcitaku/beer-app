FROM golang:1.20.1 as build

WORKDIR src/github.com/viktorcitaku/beer-app

COPY ./../../go.mod .
COPY ./../../go.sum .
COPY ./../../cmd ./cmd
COPY ./../../internal ./internal

RUN CGO_ENABLED=0 go build -ldflags "-s -w" ./cmd/beer

FROM debian:stable-slim

ENV BEER_STATIC_FILES=/opt/app

COPY ./../../web /opt/app
COPY --from=build /go/src/github.com/viktorcitaku/beer-app/beer /opt/app

CMD ["/opt/app/beer"]