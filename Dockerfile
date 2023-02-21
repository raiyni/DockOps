FROM golang:1.19-alpine as builder

WORKDIR /go/src/app

COPY go.mod ./go.mod
COPY go.sum ./go.sum

COPY cmd ./cmd
COPY pkg ./pkg

RUN apk add --update --no-cache git gcc g++ make

RUN go mod download

RUN go build -o /compose-ops cmd/main.go

# Deploy
FROM alpine:3.17 

RUN apk add --update --no-cache docker-cli

WORKDIR /

COPY --from=builder /compose-ops /compose-ops

RUN touch config.yml

ENTRYPOINT ["/compose-ops"]