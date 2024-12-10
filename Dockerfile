FROM golang:latest
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . .
CMD ["go","run","main.go"]
