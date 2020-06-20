from golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go get github.com/swaggo/swag/cmd/swag
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-X main.GIT_COMMIT=$GIT_COMMIT -X main.CONFIG_FILE=/app/config/graphSnapper-conf.json" -o main cmd/graphSnapper/main.go
RUN swag init -g cmd/graphSnapper/main.go -o ./docs

from alpine:latest
ARG CONFIG_FILE
RUN apk --no-cache add ca-certificates
RUN mkdir /app
RUN mkdir /app/config
RUN mkdir /app/logs
RUN mkdir /app/docs
COPY ${CONFIG_FILE} /app/config/graphSnapper-conf.json
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs
RUN chmod 777 /app/main
EXPOSE 8080
ENTRYPOINT ["/app/main"]
