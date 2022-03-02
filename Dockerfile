# Google App Engine の最新
FROM golang:1.16

ENV REPOSITORY=github.com/ozaki-physics/gcp-training
WORKDIR /go/src/$REPOSITORY
COPY . .
RUN go mod download
# RUN go build -o /go/bin main.go
# CMD ["/go/bin/main"]
