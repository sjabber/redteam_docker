FROM node:14-alpine3.12
MAINTAINER Tae ho Kim <sjabber91@gmail.com>

RUN mkdir -p html

WORKDIR ./html

COPY ../html .

RUN npm install

CMD ["node", "server.js"]

