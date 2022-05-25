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
go build -o bin/stacks-disposer-worker main.go
```

### Run

```shell
./bin/stacks-disposer-worker
```

### Usage

```shell
# To see current usage
./bin/stacks-disposer-worker -h
Usage of ./bin/stacks-disposer-worker:
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
IMAGE_NAME="gcr.io/superhub/stacks-disposer-worker";
IMAGE_TAG="$(git rev-parse --short HEAD)";
docker buildx build --tag "${IMAGE_NAME}:${IMAGE_TAG}" --tag "${IMAGE_NAME}:latest" . ;
docker push "${IMAGE_NAME}:${IMAGE_TAG}";
docker push "${IMAGE_NAME}:latest";
```

### Deploy service

Deploy with GCP Cloud Run

```shell
gcloud beta run deploy stacks-disposer-worker \
  --update-labels="version=$(git rev-parse --short HEAD)" \
  --format="json" \
  --region="us-central1" \
  --image="gcr.io/superhub/stacks-disposer-worker:latest" \
  --port="8080" \
  --timeout="1h" \
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
curl -i -X "DELETE" "https://stacks-disposer-worker.run.app/gke-empty-cluster/stimulating-harris-239.epam.devops.delivery?verbose=1"
```

[google-stack]: https://github.com/agilestacks/google-stacks
