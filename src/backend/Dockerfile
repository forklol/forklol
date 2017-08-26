FROM alpine
MAINTAINER Rob van Bentem <robvanbentem@gmail.com>

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

EXPOSE 8000

ADD forklol-api /

CMD ["/forklol-api", "-port", "8000"]
