# Dnote Docker Image

The official Dnote docker image.

## Installing Dnote Server Using Docker

*Installing Dnote through Docker is currently in beta.* For the an alternative installation guide, please see [the installation guide](https://github.com/dnote/dnote/blob/master/SELF_HOST.md).

### Steps

1. Install [Docker](https://docs.docker.com/install/).
2. Install Docker [Compose](https://docs.docker.com/compose/install/).
3. Download the [docker-compose.yml](https://raw.githubusercontent.com/dnote/dnote/master/host/docker/docker-compose.yml) file.

```
curl https://raw.githubusercontent.com/dnote/dnote/master/host/docker/docker-compose.yml > docker-compose.yml
```

4. Run the following to download the images and run the containers

```
docker-compose pull
docker-compose up -d
```

Visit http://localhost:3000 in your browser to see Dnote running.

please see [the installation guide](https://github.com/dnote/dnote/blob/master/SELF_HOST.md) for further configuration.
