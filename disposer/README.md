# Disposer

## Description

Web service which accept requests to run undeploy operation for specific stack from [google-stack] repo.

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

### Deploy service

Deploy with GCP Cloud Run

```shell
gcloud beta run deploy stack-disposer \
  --update-labels="version=$(git rev-parse --short HEAD)" \
  --format="json" \
  --region="us-central1" \
  --image="gcr.io/superhub/stack-disposer:latest" \
  --port="8080" \
  --allow-unauthenticated \
  --cpu-throttling \
  --execution-environment="gen2"
```

## API

This service expose next API endpoint

```api
DELETE /{sandboxId}/{stackId}
```

Where:

- `sandboxId` is a type of sandbox from [google-stack]
- `stackId` is a id of deployed stack

Also this endpoint accept `verbose` parameters to run undeploy with verbosity output.
> Note: service should be run with `-verbose` flag to see output of undeploy commands

Example of request

```shell
curl -i -X "DELETE" "https://stack-disposer.run.app/gke-empty-cluster/stimulating-harris-239.epam.devops.delivery?verbose=1"
```

[google-stack]: https://github.com/agilestacks/google-stacks
