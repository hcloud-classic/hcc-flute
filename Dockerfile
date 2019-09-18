# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /hcloud-flute/
WORKDIR /hcloud-flute/

ADD hcloud-flute /hcloud-flute/
RUN chmod 755 /hcloud-flute/hcloud-flute

EXPOSE 8001

CMD ["/hcloud-flute/hcloud-flute"]
