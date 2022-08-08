MAINTAINER vbha.mmk@gmail.com
ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}
COPY . .
RUN go build -mod=vendor -ldflags="-w -s" -o bin/controller main.go

CMD ["./bin/controller"]
