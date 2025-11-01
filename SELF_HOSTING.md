# Self-Hosting Dnote Server

Please see the [doc](https://www.getdnote.com/docs/server) for more.

## Docker Installation

1. Install [Docker](https://docs.docker.com/install/).
2. Install Docker [Compose plugin](https://docs.docker.com/compose/install/linux/).
3. Create a `compose.yml` file with the following content:

```yaml
services:
  dnote:
    image: dnote/dnote:latest
    container_name: dnote
    ports:
      - 3001:3001
    volumes:
      - ./dnote_data:/data
    restart: unless-stopped
```

4. Run the following to download the image and start the container

```
docker compose up -d
```

Visit http://localhost:3001 in your browser to see Dnote running.

## Manual Installation

Download from [releases](https://github.com/dnote/dnote/releases), extract, and run:

```bash
tar -xzf dnote-server-$version-$os.tar.gz
mv ./dnote-server /usr/local/bin
dnote-server start --baseUrl=https://your.server
```

You're up and running. Database: `~/.local/share/dnote/server.db` (customize with `--dbPath`). Run `dnote-server start --help` for options.

Set `apiEndpoint: https://your.server/api` in `~/.config/dnote/dnoterc` to connect your CLI to the server.
