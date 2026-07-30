FROM golang:1.24.4-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app/backend

COPY backend/go.mod ./
COPY backend/ ./

# go.mod currently declares the Go version but not the imported modules.
# Resolve them inside the test image without modifying the host workspace.
RUN go mod tidy

ENV CGO_ENABLED=1

CMD ["go", "test", "-v", "./..."]
