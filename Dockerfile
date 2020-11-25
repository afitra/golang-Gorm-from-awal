FROM golang:1.14.3-alpine3.11

WORKDIR /go/src/app
COPY . /go/src/app
RUN go build -o main .

ENTRYPOINT ["/app/main"]