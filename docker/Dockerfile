FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 8080

RUN make build

CMD  ./wait && ./server.out