FROM golang:1.19.4-alpine

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o ./app .

EXPOSE 8080

CMD ["./app"]