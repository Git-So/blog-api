
#build stage
FROM golang:alpine AS builder
RUN export GO111MODULE=on
RUN export GOPROXY=https://goproxy.io
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod verify
RUN go mod tidy
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ./app
LABEL Name=blog-api Version=0.0.1
EXPOSE 8099
