FROM golang:1.22-alpine AS builder

COPY . /github.com/bogdanove/auth/source/
WORKDIR /github.com/bogdanove/auth/source/

RUN go mod download
RUN go build -o ./bin/user cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/bogdanove/auth/source/bin/auth .

CMD ["./auth"]