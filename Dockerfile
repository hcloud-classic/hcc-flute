# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /GraphQL_Flute/
WORKDIR /GraphQL_Flute/

ADD GraphQL_Flute /GraphQL_Flute/
RUN chmod 755 /GraphQL_Flute/GraphQL_Flute

EXPOSE 8001

CMD ["/GraphQL_Flute/GraphQL_Flute"]
