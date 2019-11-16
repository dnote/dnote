# Dnote Docker Image

This is a Docker image to preview Dnote with one command.

## Use

```
docker run -d --name=dnote --volume=./.dnote/data:/var/lib/postgresql/data --publish 3000:3000 dnote/dnote
```
