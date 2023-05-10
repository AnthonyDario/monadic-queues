FROM golang:alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY server.go ./
COPY config.json ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /conf-server

EXPOSE 9090

CMD ["/conf-server"]
