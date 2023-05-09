FROM golang:alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY server.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /log-server

EXPOSE 8080

CMD ["/log-server"]
