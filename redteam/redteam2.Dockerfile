FROM openjdk:8-jdk
MAINTAINER Tae ho Kim <sjabber91@gmail.com>

ENV DEBIAN_FRONTEND=nointeractive
ENV TOKEN_REFRESH=awdawdawdawd
ENV TOKEN_SECRET=qlwkndlqiwndlian
ENV DB_HOST=20.194.16.227
ENV DB_USER=redteam
ENV DB_NAME=redteam
ENV DB_PW=dkagh1234!
ENV DB_PORT=5432
ENV AES_IV=0987654321654321
ENV AES_KEY=qlwkndlqiwndlian
ENV KAFKA_TOPIC=redteam
ENV KAFKA_PORT=9092
ENV KAFKA_HOSTNAME=localhost
ENV BOOT_SERVER_PORT=5001

RUN mkdir -p /api

WORKDIR /api

COPY ./mer-0.0.1-SNAPSHOT.jar .

ENTRYPOINT ["java","-jar","mer-0.0.1-SNAPSHOT.jar"] 
