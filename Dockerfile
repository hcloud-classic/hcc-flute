# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /GraphQL_flute/
WORKDIR /GraphQL_flute/

ADD GraphQL_flute /GraphQL_flute/
RUN chmod 755 /GraphQL_flute/GraphQL_flute

EXPOSE 8001

CMD ["/GraphQL_flute/GraphQL_flute"]
