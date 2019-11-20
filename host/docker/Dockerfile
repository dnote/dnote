FROM alpine:latest

ARG tarballName
RUN test -n "$tarballName"

# add dependency to execute a golang binary with dynamical linking.
RUN apk add --no-cache \
        libc6-compat

WORKDIR dnote

COPY "$tarballName" .
RUN tar -xvzf "$tarballName"

COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]

CMD ./dnote-server start

EXPOSE 3000
