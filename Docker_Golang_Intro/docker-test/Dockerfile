FROM golang:1.22-alpine
RUN mkdir /app
ENV GO111MODULE=off

WORKDIR /go/src/app

ADD . /app 
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]