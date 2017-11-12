FROM golang:1.9-alpine
MAINTAINER Trippy Salmon <cixx@cyberdyne.sh>

RUN apk --update upgrade && \
    apk add curl ca-certificates git && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/* && \
    mkdir /go/src/backend

COPY . /go/src/backend
WORKDIR /go/src/backend
RUN go get -u github.com/golang/dep/cmd/dep && \
  dep ensure && \
  go build -v -o /backend . && \
  rm -rf /go/src/backend

EXPOSE 8000

ENTRYPOINT ["/backend", "-port", "8000"]
