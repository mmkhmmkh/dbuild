FROM golang:1.19
MAINTAINER vbha.mmk@gmail.com
COPY . .
RUN go build -mod=vendor -ldflags="-w -s" -o bin/controller main.go

CMD ["./bin/controller"]
