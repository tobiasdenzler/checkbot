# Start from the latest golang base image
FROM golang:latest as builder

# Download the oc client tool
WORKDIR /tmp
ENV OC3_VERSION=v3.11.0 \
	OC3_ARCHIVE=openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit \
	OC3_SHA256SUM=4b0f07428ba854174c58d2e38287e5402964c9a9355f6c359d1242efd0990da3 \
    KUBE_VERSION=v1.17.0

RUN curl -L https://github.com/openshift/origin/releases/download/${OC3_VERSION}/${OC3_ARCHIVE}.tar.gz -o /tmp/${OC3_ARCHIVE}.tar.gz \
    && echo "${OC3_SHA256SUM} /tmp/${OC3_ARCHIVE}.tar.gz" > /tmp/${OC3_ARCHIVE}.sha256sum \
    && sha256sum -c /tmp/${OC3_ARCHIVE}.sha256sum \
    && tar xfvz /tmp/${OC3_ARCHIVE}.tar.gz --strip-components=1 -C /tmp/ \
    && rm -f /tmp/${OC3_ARCHIVE}.tar.gz

RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /tmp/kubectl \
    && chmod +x /tmp/kubectl


######## Start a new stage from scratch #######
FROM frolvlad/alpine-glibc:latest

RUN apk --no-cache update \
    && apk add --no-cache bash curl jq bind-tools python py-pip py-setuptools less coreutils \
    && apk --no-cache add ca-certificates \
    && pip --no-cache-dir install awscli \
    && rm -rf /var/cache/apk/*

RUN mkdir /app

WORKDIR /app/

# Copy the pre-built binary file
COPY main .

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
