FROM golang:1.14
MAINTAINER Tae ho Kim <sjabber91@gmail.com>

# env value
ENV AES_IV=0987654321654321
ENV AES_KEY=qlwkndlqiwndlian
ENV KAFKA_TOPIC=redteam
ENV KAFKA_PORT=9092
ENV KAFKA_HOSTNAME=localhost

RUN mkdir -p /api
WORKDIR /api

COPY producer .

ENTRYPOINT ["producer"]
