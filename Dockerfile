FROM ubuntu:22.04

RUN apt-get update -y

RUN apt install -y software-properties-common

RUN add-apt-repository ppa:longsleep/golang-backports \
    && apt update \
    && apt install -y golang-go

WORKDIR /app

COPY client client

COPY utils utils

RUN go mod init the-autoscaler && go mod tidy

WORKDIR /app/client

ARG ORCHESTRATOR_URL="localhost"

ENV ORCHESTRATOR_URL=${ORCHESTRATOR_URL}

ENTRYPOINT ["go", "run", "main.go"]