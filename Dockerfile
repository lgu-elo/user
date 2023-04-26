FROM alpine:latest

WORKDIR /app
COPY . /app

RUN apk add gcompat

EXPOSE 8082

ENTRYPOINT [ "/app/bin/user" ]
