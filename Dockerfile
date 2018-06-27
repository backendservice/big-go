FROM golang:1.10

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/lib/go/bin:$PATH

ADD . $GOPATH/src/github.com/backendservice/big-go/user-service

WORKDIR $GOPATH/src/github.com/backendservice/big-go/user-service

RUN go get -u github.com/golang/dep/...
RUN dep ensure

EXPOSE 50051
RUN go build -o server user-service/user.go
CMD ["./server"]
