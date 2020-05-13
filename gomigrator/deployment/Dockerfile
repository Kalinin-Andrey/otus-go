FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -v ./...

EXPOSE 81

CMD ["server"]
