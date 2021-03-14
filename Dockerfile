# Start from the latest golang base image
FROM ghcr.io/autamus/go:latest as builder

# Add Maintainer Info
LABEL maintainer="Alec Scott <alecbcs@github.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Add Spack Built Packages to PATH
RUN export PATH=/opt/view/bin:/opt/spack/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o buildconfig .

# Start again with minimal envoirnment.
FROM alpine

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/buildconfig /app/buildconfig

# Command to run the executable
ENTRYPOINT ["/app/buildconfig"]