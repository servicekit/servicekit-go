FROM golang:1.8-alpine

RUN apk update \
    && apk add git \
    && apk add make

RUN go get -u github.com/golang/lint/golint

WORKDIR /go/src/github.com/servicekit/servicekit-go

COPY . .

RUN mv testvendor vendor

ENTRYPOINT ["make"]