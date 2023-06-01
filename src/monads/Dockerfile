FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY server.go ./
COPY queue.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /monad-queue

EXPOSE 8901

CMD ["/monad-queue"]
