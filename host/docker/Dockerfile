FROM busybox:glibc

ARG TARGETPLATFORM
ARG version

RUN test -n "$TARGETPLATFORM" || (echo "TARGETPLATFORM is required" && exit 1)
RUN test -n "$version" || (echo "version is required" && exit 1)

WORKDIR /tmp/tarballs

# Copy all architecture tarballs
COPY dnote_server_*.tar.gz ./

# Select and extract the correct tarball based on target platform
RUN case "$TARGETPLATFORM" in \
      "linux/amd64") ARCH="amd64" ;; \
      "linux/arm64") ARCH="arm64" ;; \
      "linux/arm/v7") ARCH="arm" ;; \
      "linux/386") ARCH="386" ;; \
      *) echo "Unsupported platform: $TARGETPLATFORM" && exit 1 ;; \
    esac && \
    TARBALL="dnote_server_${version}_linux_${ARCH}.tar.gz" && \
    echo "Extracting $TARBALL for $TARGETPLATFORM" && \
    mkdir -p /dnote && \
    tar -xvzf "$TARBALL" -C /dnote

WORKDIR /dnote

# Set default database path for all processes (main server, docker exec, shells)
ENV DBPath=/data/dnote.db

COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]

CMD ./dnote-server start

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget --no-verbose --tries=1 -O /dev/null http://localhost:3001/health || exit 1
