FROM golang:1.7
RUN mkdir -p /go/src/github.com/mattmac4241/chat-auth
WORKDIR /go/src/github.com/mattmac4241/chat-auth
COPY . /go/src/github.com/mattmac4241/chat-auth

ENV PORT 8080

EXPOSE 8080

CMD ["go", "run", "main.go"]
