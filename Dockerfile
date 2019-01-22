FROM alpine:latest

MAINTAINER Chris Robertson https://github.com/electronicsleep

RUN mkdir -p /usr/src/app

RUN apk update && apk upgrade 

ADD gocode /usr/src/app
ADD public /usr/src/app/public

WORKDIR /usr/src/app
EXPOSE 8000

CMD ["./gocode"]
