FROM golang:1.19.5 as build

WORKDIR src/github.com/viktorcitaku/beer-app

COPY ./../go.mod .
COPY ./../go.sum .
COPY ./../cmd ./cmd
COPY ./../internal ./internal

RUN CGO_ENABLED=0 go build ./cmd/beer

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/src/github.com/viktorcitaku/beer-app/beer /

CMD ["/beer"]