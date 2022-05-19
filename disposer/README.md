# Disposer

## Description

Web service which accept requests to run undeploy operation for specific stack from [google-stack](https://github.com/agilestacks/google-stacks) repo.

## Development

### Prepare

```shell
go mod download
```

### Build

```shell
go build -o bin/stack-disposer main.go
```

### Run

```shell
./bin/stack-disposer
```

### Usage

```shell
# To see current usage
./bin/stack-disposer -h
Usage of ./bin/stack-disposer:
  -gitDir string
    directory where clone stacks to (default "/tmp/stacks")
  -gitUrl string
    Git URL with stacks (default "https://github.com/agilestacks/google-stacks.git")
  -port string
    port where to listen (default "8080")
  -timeout duration
    request timeout (default 1h0m0s)
  -verbose
    verbose logging
```

### Docker image

Docker image is currently based on [`gcr.io/superhub/gcp-toolbox`](https://github.com/agilestacks/toolbox/tree/master/gcp-toolbox)

To build and push image

```shell
IMAGE_NAME="gcr.io/superhub/stack-disposer";
IMAGE_TAG="$(git rev-parse --short HEAD)";
docker build --tag "${IMAGE_NAME}:${IMAGE_TAG}" --tag "${IMAGE_NAME}:latest" . ;
docker push "${IMAGE_NAME}:${IMAGE_TAG}";
docker push "${IMAGE_NAME}:latest";
```
