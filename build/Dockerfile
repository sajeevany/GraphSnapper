from golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main cmd/main.go

from alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /app/main /app/main
RUN chmod 777 /app/main
EXPOSE 8080
ENTRYPOINT ["/app/main"]
