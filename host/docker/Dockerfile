FROM postgres:11.5

ARG version
RUN test -n "$version"

ENV POSTGRES_USER=dnote_preview
ENV POSTGRES_PASSWORD=dnote_preview
ENV POSTGRES_DB=dnote_preview

WORKDIR dnote

RUN apt-get update
COPY dnote_server_"$version"_linux_amd64.tar.gz .
RUN tar -xvzf dnote_server_"$version"_linux_amd64.tar.gz

ADD docker-entry.sh .
RUN chmod +x ./docker-entry.sh
ENTRYPOINT ./docker-entry.sh

EXPOSE 3000
