# Self-Hosting Dnote Server

Please see the [doc](https://www.getdnote.com/docs/server) for more.

## Docker Installation

1. Install [Docker](https://docs.docker.com/install/).
2. Install Docker [Compose plugin](https://docs.docker.com/compose/install/linux/).
3. Download the [compose.yml](https://raw.githubusercontent.com/dnote/dnote/master/host/docker/compose.yml) file by running:

```
curl https://raw.githubusercontent.com/dnote/dnote/master/host/docker/compose.yml > compose.yml
```

4. Run the following to download the images and run the containers

```
docker compose pull
docker compose up -d
```

Visit http://localhost:3001 in your browser to see Dnote running.

### Supported platform

Currently, the official Docker image for Dnote supports Linux running AMD64 CPU architecture.

If you run ARM64, please install Dnote server by downloading a binary distribution (see below).

## Manual Installation

Download from [releases](https://github.com/dnote/dnote/releases), extract, and run:

```bash
tar -xzf dnote-server-$version-$os.tar.gz
mv ./dnote-server /usr/local/bin
dnote-server start --webUrl=https://your.server
```

You're up and running. Database: `~/.local/share/dnote/server.db` (customize with `--dbPath`). Run `dnote-server start --help` for options.

Set `apiEndpoint: https://your.server/api` in `~/.config/dnote/dnoterc` to connect your CLI to the server.
