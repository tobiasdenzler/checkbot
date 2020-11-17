# Start from the latest golang base image
FROM golang:latest as builder

# Some build arguments
ARG GIT_VERSION=unspecified
LABEL git_version=$GIT_VERSION
ARG GIT_BUILD=unspecified
LABEL git_build=$GIT_BUILD

# Download the oc client tools
WORKDIR /tmp
ENV OC3_VERSION=v3.11.0 \
    OC3_ARCHIVE=openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit \
    OC3_SHA256SUM=4b0f07428ba854174c58d2e38287e5402964c9a9355f6c359d1242efd0990da3 \
    OC4_VERSION=4.5.0-0.okd-2020-10-15-235428 \
    OC4_SHA256SUM=cdb5cd07eb18c7dbeff9e817d2221a42d29cfd0f4c129f41b519ce06772d54c0 \
    KUBE_VERSION=v1.19.0 \
    AWS_VERSION=2.1.0

# OCP3
RUN curl -L https://github.com/openshift/origin/releases/download/${OC3_VERSION}/${OC3_ARCHIVE}.tar.gz -o /tmp/${OC3_ARCHIVE}.tar.gz \
    && echo "${OC3_SHA256SUM} /tmp/${OC3_ARCHIVE}.tar.gz" > /tmp/${OC3_ARCHIVE}.sha256sum \
    && sha256sum -c /tmp/${OC3_ARCHIVE}.sha256sum \
    && tar xfvz /tmp/${OC3_ARCHIVE}.tar.gz --strip-components=1 -C /tmp/ \
    && mv /tmp/oc /tmp/oc3 \
    && rm -f /tmp/${OC3_ARCHIVE}.tar.gz

# OCP4
RUN curl -L https://github.com/openshift/okd/releases/download/${OC4_VERSION}/openshift-client-linux-${OC4_VERSION}.tar.gz -o /tmp/openshift-client-linux-${OC4_VERSION}.tar.gz \
    && echo "${OC4_SHA256SUM} /tmp/openshift-client-linux-${OC4_VERSION}.tar.gz" > /tmp/openshift-client-linux-${OC4_VERSION}.sha256sum \
    && sha256sum -c /tmp/openshift-client-linux-${OC4_VERSION}.sha256sum \
    && tar xfvz /tmp/openshift-client-linux-${OC4_VERSION}.tar.gz -C /tmp/ \
    && mv /tmp/oc /tmp/oc4 \
    && rm -f /tmp/openshift-client-linux-${OC4_VERSION}.tar.gz

# K8S
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl -o /tmp/kubectl \
    && chmod +x /tmp/kubectl

# AWS
RUN curl -L https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${AWS_VERSION}.zip -o /tmp/awscliv2.zip \
    && apt-get update \
    && apt-get -y install unzip \
    && unzip -q /tmp/awscliv2.zip

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-w -s -X main.Version=${GIT_VERSION} -X main.Build=${GIT_BUILD}" -a -installsuffix cgo -o main ./cmd/server/


######## Start a new stage from scratch #######
FROM frolvlad/alpine-glibc:latest

RUN apk --no-cache update \
    && apk add --no-cache bash curl jq bind-tools python3 py-pip py-setuptools less coreutils \
    && apk --no-cache add ca-certificates \
    && pip --no-cache-dir install awscli \
    && rm -rf /var/cache/apk/*

RUN mkdir /app

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Add the oc3 client tool
COPY --from=builder /tmp/oc3 /usr/bin/

# Provide backwards compatibility for scripts
RUN ln -s /usr/bin/oc3 /usr/bin/oc

# Add the oc4 client tool
COPY --from=builder /tmp/oc4 /usr/bin/

# Add the kubectl client tool
COPY --from=builder /tmp/kubectl /usr/bin/

# Copy the AWS CLI
COPY --from=builder /tmp/aws/ /tmp/aws/

# Install the AWS CLI
RUN /tmp/aws/install

# Copy the certs
COPY certs certs

# Copy static ui files
COPY ui ui

# Expose port to the outside world
EXPOSE 4444

# Command to run the executable
CMD ["./main"] 
