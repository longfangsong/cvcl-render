
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cvcl-render

# Download modern-cv package
RUN apk add --no-cache git && \
    mkdir -p /builder-typst-packages && \
    cd /builder-typst-packages && \
    git clone --depth 1 https://github.com/longfangsong/modern-cv.git && \
    ls -la modern-cv/

# Runtime stage
FROM ghcr.io/typst/typst:latest

# Install necessary fonts: Roboto, Source Sans Pro, and FontAwesome
RUN apk add --no-cache \
    font-noto \
    font-noto-cjk \
    font-noto-extra \
    ttf-dejavu \
    ttf-freefont \
    ttf-liberation \
    ttf-droid \
    unzip && \
    mkdir -p /usr/share/fonts/ttf-roboto /usr/share/fonts/ttf-source-sans-pro /usr/share/fonts/ttf-font-awesome && \
    wget -qO- https://github.com/google/fonts/raw/refs/heads/main/ofl/roboto/Roboto%5Bwdth,wght%5D.ttf > /usr/share/fonts/ttf-roboto/Roboto.ttf && \
    wget -qO /tmp/fontawesome.zip https://use.fontawesome.com/releases/v7.1.0/fontawesome-free-7.1.0-desktop.zip && \
    unzip -q /tmp/fontawesome.zip -d /tmp && \
    cp -r /tmp/fontawesome-free-7.1.0-desktop/otfs/*.otf /usr/share/fonts/ttf-font-awesome/ && \
    rm -rf /tmp/fontawesome.zip /tmp/fontawesome-free-7.1.0-desktop
COPY ./font/ /usr/share/fonts/ttf-source-sans-pro/
RUN fc-cache -fv && mkdir -p /root/.local/share/typst/packages/local/modern-cv/0.9.0

# Copy the built binary from builder stage
COPY --from=builder /build/cvcl-render /usr/local/bin/cvcl-render
# Copy modern-cv package from builder
COPY --from=builder /builder-typst-packages/modern-cv /root/.local/share/typst/packages/local/modern-cv/0.9.0/

# Create output directory
RUN mkdir -p /output

# Set working directory
WORKDIR /app

# Expose default port
EXPOSE 8080

# Set entrypoint to cvcl-render
ENTRYPOINT ["cvcl-render"]

# Default command - run HTTP server
CMD ["-output-dir", "/output"]
