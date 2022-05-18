# Stack disposer

## Description

This repo contains all required services needed to implement auto cleanup of deployed stacks in Google Cloud Provider

There are next services:

* Disposer - service provides endpoint to run undeploy operation on given stackId
* Scanner - GCP Cloud Function which monitor deployed stacks and send undeploy requests to Disposer based on time criteria (ex. stack unused for 7 days) WIP
