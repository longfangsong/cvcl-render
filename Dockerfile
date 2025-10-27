
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cvcl-render

# Runtime stage
FROM ghcr.io/typst/typst:latest

# Install necessary fonts: Roboto, Source Sans Pro, and FontAwesome
RUN apk add --no-cache \
    font-noto \
    font-noto-cjk \
    font-noto-extra \
    ttf-font-awesome \
    ttf-dejavu \
    ttf-freefont \
    ttf-liberation \
    ttf-droid && \
    mkdir -p /usr/share/fonts/ttf-roboto /usr/share/fonts/ttf-source-sans-pro && \
    wget -qO- https://github.com/google/fonts/raw/refs/heads/main/ofl/roboto/Roboto%5Bwdth,wght%5D.ttf > /usr/share/fonts/ttf-roboto/Roboto.ttf
COPY ./SourceSansPro-Regular.otf /usr/share/fonts/ttf-source-sans-pro/SourceSansPro-Regular.otf
RUN fc-cache -fv

# Copy the built binary from builder stage
COPY --from=builder /build/cvcl-render /usr/local/bin/cvcl-render

# Create output directory
RUN mkdir -p /output

# Set working directory
WORKDIR /app

# Expose default port
EXPOSE 8080

# Default command - run HTTP server
CMD ["cvcl-render", "-output-dir", "/output"]
