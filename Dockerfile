
#build stage
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN export GO111MODULE=on
RUN if [ ! -d "vendor" ];then go mod vendor;fi
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -installsuffix cgo . 

#final stage
FROM alpine:latest
COPY --from=builder /app/blog-api /blog-api
ENTRYPOINT /blog-api
LABEL Name=blog-api Version=0.0.1

