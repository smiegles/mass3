# Build Container
FROM golang:1.10.3-alpine3.7 AS build-env
RUN apk add --no-cache --upgrade git openssh-client ca-certificates
RUN go get -u github.com/golang/dep/cmd/dep

# Need to use full path to handle local package imports (argo and common)
WORKDIR /go/src/github.com/anshumanbh/mass3

# Cache the dependencies early
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only -v

# Build
COPY mass3.go ./
RUN go build -v -o ${GOPATH}/bin/mass3

# Final Container
FROM alpine:3.7
LABEL maintainer="Anshuman Bhartiya"
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/mass3 /usr/bin/mass3
COPY lists/resolvers.txt ./
COPY lists/buckets.txt ./
RUN apk update && apk add --no-cache --upgrade bash bind-tools

ENTRYPOINT ["/usr/bin/mass3"]

