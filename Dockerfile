FROM golang as builder

ARG VERSION

WORKDIR /app

COPY . /app/

RUN apt update && \
    apt install npm -y 
RUN make install-js && \
    make version=${VERSION} build-web

RUN go get -u github.com/gobuffalo/packr/v2/packr2 && \
    make version=${VERSION} build-server && \
    tar -xzvf /app/build/server/dnote_server_${VERSION}_linux_amd64.tar.gz


FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/dnote-server .

ENTRYPOINT ["/app/dnote-server", "start"]
