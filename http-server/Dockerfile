FROM golang:1.21.1

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o parkpal-backend cmd/main.go

EXPOSE 9091

CMD ["./parkpal-backend"]