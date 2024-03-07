FROM golang:1.22.0-alpine3.18

LABEL maintainer="Azraf Al Monzim"
LABEL name="Test" version="1.0.0"
LABEL author="Azraf Al Monzim"

WORKDIR /go/src/app

COPY . .

# create a .env file
RUN echo "PORT=8080" >.env

RUN go get -d -v

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
