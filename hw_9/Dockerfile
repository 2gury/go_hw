FROM golang:1.17-alpine AS builder

COPY . /go_ci-cd/
WORKDIR /go_ci-cd/

RUN go mod download
RUN go build -o ./.bin/hello cmd/main.go

FROM alpine:latest

COPY --from=builder /go_ci-cd/.bin/hello .

EXPOSE 8080

CMD ["./hello"]

