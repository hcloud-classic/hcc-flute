# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /flute/
WORKDIR /flute/

ADD flute /flute/
RUN chmod 755 /flute/flute

EXPOSE 7000

CMD ["/flute/flute"]
