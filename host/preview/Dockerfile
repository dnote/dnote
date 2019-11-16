FROM postgres:11.5

ENV DNOTE_SERVER_VERSION=0.3.1
ENV POSTGRES_USER=dnote_preview
ENV POSTGRES_PASSWORD=dnote_preview
ENV POSTGRES_DB=dnote_preview

WORKDIR dnote

RUN apt-get update
ADD https://github.com/dnote/dnote/releases/download/server-v"$DNOTE_SERVER_VERSION"/dnote_server_"$DNOTE_SERVER_VERSION"_linux_amd64.tar.gz .
RUN tar -xvzf dnote_server_"$DNOTE_SERVER_VERSION"_linux_amd64.tar.gz

ADD docker-entry.sh .
RUN chmod +x ./docker-entry.sh
ENTRYPOINT ./docker-entry.sh

EXPOSE 3000
