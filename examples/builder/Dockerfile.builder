FROM debian:jessie

RUN apt-get update && \
    apt-get install -y golang && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Create src directory for compiling artifacts
ADD . /src
WORKDIR /src

# Get all package dependencies
RUN go get -t -d -v ./...
# Build go binary
RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags=--s -o hello

# Create directory for final image
RUN mkdir /build

# Copy runner Dockerfile and hello binary to build dir
RUN cp Dockerfile.runner /build/Dockerfile
RUN cp hello /build/

# Change workdir to build
WORKDIR "/build"

# When builder runs stream a tar file of the build directory
CMD tar -cf - .
