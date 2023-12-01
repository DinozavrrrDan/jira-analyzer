FROM ubuntu:latest
LABEL authors="Boris"

ENTRYPOINT ["top", "-b"]