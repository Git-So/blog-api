#final stage
FROM alpine:latest
COPY /app/blog-api /blog-api
COPY  .blog .blog
ENTRYPOINT ./blog-api
LABEL Name=blog-api Version=0.0.1
EXPOSE 8099
