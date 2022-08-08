FROM golang:1.19
MAINTAINER vbha.mmk@gmail.com
WORKDIR /dbuild
COPY . .
RUN go build -mod=vendor -ldflags="-w -s" -o bin/controller controller/main.go

CMD ["/dbuild/bin/controller"]
