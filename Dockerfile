# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
ADD . /go/src/github.com/Andrew-Klaas/go-movies-app
WORKDIR /go/src/github.com/Andrew-Klaas/go-movies-app
RUN go get github.com/satori/go.uuid
RUN go get github.com/hashicorp/vault/api
RUN go get github.com/lib/pq
RUN go install /go/src/github.com/Andrew-Klaas/go-movies-app

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/go-movies-app

# Document that the service listens on port 8080.
EXPOSE 8080

#docker build -t aklaas2/go-movies-app .;docker push aklaas2/go-movies-app:latest
#docker build -t aklaas2/go-movies-app-v2 .;docker push aklaas2/go-movies-app-v2:latest