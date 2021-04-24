FROM golang:1.12.7-alpine3.10 AS build
WORKDIR /go/src/app
COPY ./ ./
RUN go mod download
RUN GOOS=linux go build -o ./bin/app ./cmd/main.go

FROM alpine:3.10
WORKDIR /usr/local/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 8080
ENTRYPOINT /go/bin/app
