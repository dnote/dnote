# Dnote Docker Image

The official Dnote docker image.

## Installing Dnote Server Using Docker

1. Install [Docker](https://docs.docker.com/install/).
2. Install Docker [Compose](https://docs.docker.com/compose/install/).
3. Download the [docker-compose.yml](https://raw.githubusercontent.com/dnote/dnote/master/host/docker/docker-compose.yml) file by running:

```
curl https://raw.githubusercontent.com/dnote/dnote/master/host/docker/docker-compose.yml > docker-compose.yml
```

4. Run the following to download the images and run the containers

```
docker-compose pull
docker-compose up -d
```

Visit http://localhost:3000 in your browser to see Dnote running.

Please see [the installation guide](https://github.com/dnote/dnote/blob/master/SELF_HOSTING.md) for further configuration.

## Supported platform

Currently, the official Docker image for Dnote supports Linux running AMD64 CPU architecture.

If you run ARM64, please install Dnote server by downloading a binary distribution by following [the guide](https://github.com/dnote/dnote/blob/master/SELF_HOSTING.md).
