
#build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN export GO111MODULE=on
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ./app
LABEL Name=blog-api Version=0.0.1

