FROM golang:1.22-alpine

# Add postgresql-client for database check
RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 3120

# Create entrypoint script with proper environment variables
RUN echo "#!/bin/sh" > /entrypoint.sh && \
    echo "export DB_USER=root" >> /entrypoint.sh && \
    echo "export DB_PASSWORD=P@ssw0rd" >> /entrypoint.sh && \
    echo "until pg_isready -h db -U root -d user_service; do" >> /entrypoint.sh && \
    echo "  echo 'Waiting for postgres...';" >> /entrypoint.sh && \
    echo "  sleep 2;" >> /entrypoint.sh && \
    echo "done" >> /entrypoint.sh && \
    echo "./main" >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]