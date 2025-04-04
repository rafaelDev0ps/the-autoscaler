FROM ubuntu:22.04

RUN apt-get update -y

WORKDIR /app

COPY ../client client
COPY ../utils utils

RUN go mod tidy

EXPOSE 8081

ENTRYPOINT ["tail", "-f", "/var/log/null"]