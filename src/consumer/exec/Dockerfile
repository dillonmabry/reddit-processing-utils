FROM golang:1.9

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go get -d -v ./...
RUN go build -o main .

CMD ["/app/main", "--queue", "test"]
