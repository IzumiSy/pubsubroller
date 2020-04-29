# pubsubroller
[![CircleCI](https://circleci.com/gh/IzumiSy/pubsubroller.svg?style=svg)](https://circleci.com/gh/IzumiSy/pubsubroller)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

GCP PubSub provisioning tool at a light speed :zap:

## Table of Contents
- [Install](#Install)
- [Usage](#Usage)
- [Troubleshooting](#Troubleshooting)
- [Development](#Development)
- [Contributions](#Contributions)
- [License](#License)

## Install
```bash
$ go get github.com/IzumiSy/pubsubroller
```

## Usage
```bash
 $ pubsubroller --help
Usage:
  pubsubroller [OPTIONS]

Application Options:
  -p, --projectId= target GCP project ID
  -c, --config=    configuration file path
  -e, --endpoint=  service endpoint
      --dry        dry run
      --delete     delete all topics and their subscriptions

Help Options:
  -h, --help       Show this help messag
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

## Contributions
PRs accepted

## License
MIT © IzumiSy
