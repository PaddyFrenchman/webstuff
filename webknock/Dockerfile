FROM debian:bookworm as BUILD
# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    GO111MODULE=on
# Install dependencies
RUN apt-get update 
RUN apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        git \
        build-essential \
        golang \
        wget
RUN rm -rf /var/lib/apt/lists/*

# Install Go
RUN export GO_VERSION="1.24.5"
# RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go.tar.gz 
WORKDIR /tmp
# RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
# RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go${GO_VERSION}.linux-amd64.tar.gz
RUN curl -fsSL https://go.dev/dl/go1.24.5.linux-amd64.tar.gz -o /tmp/go${GO_VERSION}.linux-amd64.tar.gz
RUN rm -rf /usr/local/go
RUN tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz
RUN rm /tmp/go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH=/usr/local/go/bin:$PATH

# Verify installation
RUN go version

WORKDIR /app
# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
# RUN go mod download
# Copy the source code
COPY main.go ./
# Copy the self-signed certificates
RUN go build

#CMD [ "bash" ]


# Use official Golang image as the base
FROM debian:bookworm as WEBKNOCK

# Install SQLite dependencies
RUN apt-get update && apt-get upgrade
RUN apt-get install -y sqlite3
RUN apt-get install -y libsqlite3-0 
RUN apt-get install -y libsqlite3-dev 
#RUN apt-get install -y sqlite-dev 
RUN apt-get install -y gcc 
RUN apt-get install -y musl-dev

RUN mkdir /app
# Set working directory
WORKDIR /app
COPY static/ ./static/

#copies the build
COPY --from=BUILD /app/webknock ./

# Copy the self-signed certificates
COPY cert.pem key.pem ./

# Expose port 8080 for HTTPS
EXPOSE 8080

# Command to run the application
CMD ["/app/webknock"]