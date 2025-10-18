FROM busybox:glibc

ARG tarballName
RUN test -n "$tarballName"

WORKDIR dnote

COPY "$tarballName" .
RUN tar -xvzf "$tarballName"

COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]

CMD ./dnote-server start

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget --no-verbose --tries=1 -O /dev/null http://localhost:3001/health || exit 1
