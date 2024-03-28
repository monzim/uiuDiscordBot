# Stage 1: Build the Go application
FROM golang:1.22.0-alpine3.18 AS builder

WORKDIR /go/src/app

COPY . .

RUN go get -d -v
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19.1

LABEL name="UIU Discord Bot" version="1.0.0"
LABEL author="Azraf Al Monzim"
LABEL maintainer="Azraf Al Monzim"

WORKDIR /app

COPY --from=builder /go/src/app/main .

RUN echo "PORT=8080" >.env

EXPOSE 8080

CMD ["./main"]
