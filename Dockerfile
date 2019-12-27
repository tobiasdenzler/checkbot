# Start from the latest golang base image
FROM golang:latest as builder

# Download the oc client tool
WORKDIR /tmp
ENV OC3_VERSION=v3.11.0 \
	OC3_ARCHIVE=openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit \
	OC3_SHA256SUM=4b0f07428ba854174c58d2e38287e5402964c9a9355f6c359d1242efd0990da3 \
    KUBE_VERSION=v1.17.0

ADD https://github.com/openshift/origin/releases/download/${OC3_VERSION}/${OC3_ARCHIVE}.tar.gz .
RUN echo "${OC3_SHA256SUM} /tmp/${OC3_ARCHIVE}.tar.gz" > /tmp/${OC3_ARCHIVE}.sha256sum \
    && sha256sum -c /tmp/${OC3_ARCHIVE}.sha256sum \
    && tar xfvz /tmp/${OC3_ARCHIVE}.tar.gz --strip-components=1 -C /tmp/ \
    && rm -f /tmp/${OC3_ARCHIVE}.tar.gz

RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /tmp/kubectl \
    && chmod +x /tmp/kubectl

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/


######## Start a new stage from scratch #######
FROM frolvlad/alpine-glibc:latest

RUN apk --no-cache update \
    && apk add --no-cache bash curl jq bind-tools python py-pip py-setuptools less coreutils \
    && apk --no-cache add ca-certificates \
    && pip --no-cache-dir install awscli \
    && apk del --purge deps \
    && rm -rf /var/cache/apk/*

RUN mkdir /app

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Add the oc client tool
COPY --from=builder /tmp/oc /usr/bin/

# Add the kubectl client tool
COPY --from=builder /tmp/kubectl /usr/bin/

# Copy the certs
COPY certs certs

# Copy static ui files
COPY ui ui

# Expose port to the outside world
EXPOSE 4444

# Command to run the executable
CMD ["./main"] 
