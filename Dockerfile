FROM golang:1.24

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y wget
RUN mkdir -p /scripts && \
    wget -qO- https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh > /scripts/wait-for-it.sh && \
    chmod +x /scripts/wait-for-it.sh

EXPOSE 8080

CMD ["sh", "-c", "/scripts/wait-for-it.sh postgres:5432 --timeout=60 && /scripts/wait-for-it.sh redis:6379 --timeout=60 && air"]