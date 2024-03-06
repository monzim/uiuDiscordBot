FROM golang:1.22.0-alpine3.18

LABEL maintainer="Azraf Al Monzim"
LABEL name="test" version="1.0.0"
LABEL author="Azraf Al Monzim"

WORKDIR /app

RUN touch .env

COPY . /app

RUN go get -d -v

RUN go build -o test

CMD ["./test"]
