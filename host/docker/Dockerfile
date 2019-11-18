FROM alpine:latest

ARG version
RUN test -n "$version"

# add dependency to execute a golang binary with dynamical linking.
RUN apk add --no-cache \
        libc6-compat

WORKDIR dnote

COPY dnote_server_"$version"_linux_amd64.tar.gz .
RUN tar -xvzf dnote_server_"$version"_linux_amd64.tar.gz

COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]

CMD ./dnote-server start

EXPOSE 3000
