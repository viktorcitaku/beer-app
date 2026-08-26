FROM golang:1.27.0 AS build

WORKDIR /src/github.com/viktorcitaku/beer-app

COPY ./../../go.mod .
COPY ./../../go.sum .
COPY ./../../cmd ./cmd
COPY ./../../internal ./internal

RUN go mod download && go mod verify
RUN CGO_ENABLED=0 go build -ldflags "-s -w" ./cmd/beer

FROM debian:stable-slim

ENV BEER_STATIC_FILES=/opt/app

COPY ./../../web /opt/app
COPY --from=build /src/github.com/viktorcitaku/beer-app/beer /opt/app

CMD ["/opt/app/beer"]
