FROM golang:1.17.6-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY internal ./internal

RUN go build -o /main

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /main /main

COPY db/migrations /db
COPY configs /configs

EXPOSE 9090

ENTRYPOINT ["/main"]
