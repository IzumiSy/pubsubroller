# pubsubroller

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
  -endpoint string
    	service endpoint # you can use this for pub/sub emulator
  -projectId string
    	target GCP project ID
```

## Configuration Example
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
```
