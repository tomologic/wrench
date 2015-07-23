FROM golang

ADD . /src
WORKDIR /src

# Get all package dependencies
RUN mkdir /build
RUN go get -t -d -v ./...
RUN CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags=--s -o /build/hello
RUN cp Dockerfile.runner /build/Dockerfile

WORKDIR "/build"

CMD tar -cf - .
