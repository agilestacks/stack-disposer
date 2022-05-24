# Scanner

## Description

GCP Cloud Function which scan stacks and send request to disposer service to undeploy them

## Development

### Prepare

Cloud Function created using Node.js runtime. So you should have installed Node.js with version 16+.

To install dependencies

```shell
npm install
```

### Run locally

`@google-cloud/functions-framework` is used to run this function locally. It wrap the function in express.js framework and start the server

```shell
npm start
```

### Debug locally

It's also posible to debug function locally by running

```shell
npm run debug
```

And after that attach debugger to it. Example of launch config for VS Code

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Attach",
            "port": 9229,
            "request": "attach",
            "skipFiles": [
                "<node_internals>/**"
            ],
            "type": "node"
        }
    ]
}
```

### Deploy

To deploy function to GCP you should have configured `gcloud` and have appropriate permissions to work with Cloud Function

```shell
gcloud functions deploy "stack-disposer-scanner" \
    --format=json \
    --runtime=nodejs16 \
    --trigger-http \
    --allow-unauthenticated \
    --source="." \
    --entry-point="scan" \
    --set-env-vars='STATE_FUNCTION_URL=https://us-central1-superhub.cloudfunctions.net/stacks,DISPOSER_URL=https://stack-disposer-mvn4dxj74a-uc.a.run.app,DAYS_BEFORE=7,VERBOSE=false,TARGET_STATUSES=deployed;incomplete'
```

Also, you need to create a Cloud Scheduler to invoke function be schedule.

```shell
gcloud scheduler jobs create http "stacks-disposer-job" \
    --format="json" \
    --location="us-central1" \
    --schedule="0 0 * * *" \
    --uri="https://us-central1-superhub.cloudfunctions.net/stack-disposer-scanner" \
    --http-method="GET"
```
