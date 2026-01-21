# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Build stage for Go backend
FROM golang:1.24-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=1 GOOS=linux go build -o cms-server .

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=backend-builder /app/cms-server .

RUN mkdir -p /app/uploads /app/data

ENV GIN_MODE=release
ENV PORT=8080
ENV DB_PATH=/app/data/cms.db
ENV UPLOAD_DIR=/app/uploads

EXPOSE 8080

VOLUME ["/app/data", "/app/uploads"]

CMD ["./cms-server"]
