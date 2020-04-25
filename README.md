# pubsubroller
[![CircleCI](https://circleci.com/gh/IzumiSy/pubsubroller.svg?style=svg)](https://circleci.com/gh/IzumiSy/pubsubroller)

GCP PubSub provisioning tool at a light speed :zap:

## Install
```bash
$ go get github.com/IzumiSy/pubsubroller
```

## Usage
```bash
 $ pubsubroller -help
Usage of pubsubroller:
  -config string
    	configuration file path
  -delete
    	delete all topics and their subscriptions
  -dry
    	dry run
  -endpoint string
    	service endpoint
```

### Configuration Example
```yaml
variables:
  url: "https://service-of-${projectId}/subscriber"
topics:
  invitedUser:
    subscriptions:
      - name: sendInvitationMail
        endpoint: "${url}/sendInvitationMail"
      - name: sendGroupNotification
        endpoint: "${url}/sendGroupNotification"
      - name: sendReinvitationMail
        pull: true
```
top level keys are `variables`, which replaces placeholder in subscription names and endpoints, and `topics` that have multiple subscriptions.

## Troubleshooting

### `panic: rpc error: code = Unauthenticated desc = transport: oauth2: cannot fetch token: 400 Bad Request`
Try checking if your Google Cloud credential is valid or not.
```bash
$ gcloud auth application-default print-access-token
```
Revoke it and re-login if it is not valid.
```bash
$ gcloud auth application-default revoke
$ gcloud auth application-default login
```

## Development
```bash
$ make
```
