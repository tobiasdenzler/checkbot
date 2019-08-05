# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download the oc client tool
ADD https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz .
RUN tar xfvz openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz \
    && rm -f openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

######## Start a new stage from scratch #######
FROM frolvlad/alpine-glibc:latest

RUN apk add --no-cache bash curl \
    && apk --no-cache add ca-certificates

RUN mkdir /app

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Add the oc client tool
COPY --from=builder /app/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit/oc /usr/bin/

# Copy scripts and change permissions
COPY --from=builder /app/scripts ./scripts
RUN chmod -R 777 ./scripts

# Expose port to the outside world
EXPOSE 4444

# Command to run the executable
CMD ["./main"] 
