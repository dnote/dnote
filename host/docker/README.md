# Dnote Docker Image

This is a Docker image to preview Dnote with one command.

## Use

```
docker run --name=dnote_preview --volume=/.dnote/server/data:/var/lib/postgresql/data -d --publish 3000:3000 dnote/dnote
```
