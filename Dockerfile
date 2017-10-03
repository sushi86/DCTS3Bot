FROM golang:latest
RUN mkdir /app
ADD . /go/src/github.com/sushi86/DCTS3Bot
WORKDIR /go/src/github.com/sushi86/DCTS3Bot
RUN go get ./...
RUN go build -o bot .
CMD ["/go/src/github.com/sushi86/DCTS3Bot/bot"]