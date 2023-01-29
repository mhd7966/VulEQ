FROM golang:1.16-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ARG BUILD_PATH=./cmd/api
RUN apk add build-base
RUN go install github.com/rubenv/sql-migrate/sql-migrate@latest
COPY . .
RUN go build -ldflags="-w -s" -o go-app $BUILD_PATH

FROM sonarsource/sonar-scanner-cli:4
# FROM debian:11.0
# RUN apt update && apt install -y wget unzip
# RUN wget https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.6.2.2472-linux.zip && unzip sonar-scanner-cli-4.6.2.2472-linux.zip -d /sonarscanner && rm sonar-scanner-cli-4.6.2.2472-linux.zip
COPY --from=builder /app/go-app /app/go-app
COPY --from=builder /go/bin/sql-migrate /usr/local/bin/sql-migrate
COPY ./pkg/migrations/ /pkg/migrations/
COPY dbconfig.yml .
CMD ["/app/go-app"]
