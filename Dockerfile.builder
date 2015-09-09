FROM debian:jessie

RUN echo "deb http://http.debian.net/debian jessie-backports main" >> /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y git golang gcc && \
    apt-get -y -t jessie-backports install "docker.io" && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Create src directory for compiling artifacts
ENV GOPATH=/go
ADD . /go/src/github.com/tomologic/wrench
WORKDIR /go/src/github.com/tomologic/wrench

# Get all package dependencies
RUN go get -t -d -v ./...
# Build go binary
RUN go build -a -ldflags "-w -X github.com/tomologic/wrench/version.version '$(git describe)'" -o wrench

# Create directory for final image
RUN mkdir /build

# Copy runner Dockerfile and wrench binary to build dir
RUN cp Dockerfile.runner /build/Dockerfile
RUN cp wrench /build/

# Change workdir to build
WORKDIR "/build"

# When builder runs stream a tar file of the build directory
CMD tar -cf - .
