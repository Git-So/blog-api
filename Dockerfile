
#build stage
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN export GO111MODULE=on
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/blog-api /blog-api
ENTRYPOINT ./blog-api
LABEL Name=blog-api Version=0.0.1

