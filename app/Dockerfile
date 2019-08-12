FROM golang:latest

LABEL maintainer="SHREY SONI <shreysoni1@gmail.com>"

WORKDIR /block

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o block

EXPOSE 8080

CMD ["./block"]
