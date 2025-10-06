FROM busybox:glibc

ARG tarballName
RUN test -n "$tarballName"

WORKDIR dnote

COPY "$tarballName" .
RUN tar -xvzf "$tarballName"

COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]

CMD ./dnote-server start

EXPOSE 3000
