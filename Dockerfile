# syntax=docker/dockerfile:1

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.16-alpine
WORKDIR /go/src/github.com/Andrew-Klaas/go-movies-app
ADD . /go/src/github.com/Andrew-Klaas/go-movies-app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go get github.com/Andrew-Klaas/go-movies-app
RUN go install /go/src/github.com/Andrew-Klaas/go-movies-app

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/go-movies-app

# Document that the service listens on port 8080.
EXPOSE 8080

#docker build -t aklaas2/go-movies-app .;docker push aklaas2/go-movies-app:latest
#docker build -t aklaas2/go-movies-app-v2 .;docker push aklaas2/go-movies-app-v2:latest
