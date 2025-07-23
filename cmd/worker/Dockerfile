# Stage 1: Python + dependencies
FROM python:3.11-slim AS python-base

WORKDIR /app

# Install system deps and requirements
RUN apt-get update && apt-get install -y git build-essential libgl1 ffmpeg && rm -rf /var/lib/apt/lists/*
COPY scripts/requirements.txt /app/
RUN pip install --upgrade pip && pip install -r requirements.txt

# Stage 2: Golang build
FROM golang:1.24-alpine AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o worker ./cmd/worker

# Stage 3: Final combined image
FROM python:3.11-slim

WORKDIR /app

# Copy pre-built Python env
COPY --from=python-base /usr/local/lib/python3.11 /usr/local/lib/python3.11
COPY --from=python-base /usr/local/bin /usr/local/bin

# Copy scripts separately (this is the fast-changing layer)
COPY scripts/ /app/scripts/

# Copy Go binary
COPY --from=go-builder /app/worker /app/

CMD ["./worker"]
