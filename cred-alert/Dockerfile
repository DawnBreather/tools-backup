FROM golang:alpine as gobuild

RUN apk --no-cache add unzip
RUN apk add --no-cache --virtual .build-dependancies git \
    && go get github.com/pivotal-cf/cred-alert \
    && go install github.com/pivotal-cf/cred-alert \
    && apk del .build-dependancies


FROM python:3.6.0-alpine

RUN mkdir /etc/cred-alert
COPY ./credentials.regexp /etc/cred-alert/regexp
RUN apk --no-cache add unzip git

COPY --from=gobuild /go/bin/cred-alert /usr/bin
RUN pip install detect-secrets

RUN mkdir -p /home/secrets-audit
WORKDIR /home/secrets-audit

ENV PATH "$PATH:/etc/cred-alert"

