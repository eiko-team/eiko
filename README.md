# eiko
![maintenance](https://img.shields.io/maintenance/yes/2019)
[![CircleCI](https://circleci.com/gh/eiko-team/eiko.svg?style=svg)](https://app.circleci.com/github/eiko-team/eiko/pipelines)
[![Go Report Card](https://goreportcard.com/badge/github.com/eiko-team/eiko)](https://goreportcard.com/report/github.com/eiko-team/eiko)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/13cbb61d7e734f16a8f0494e0a13a993)](https://www.codacy.com/manual/tomMoulard/eiko?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=eiko-team/eiko&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/eiko-team/eiko/branch/master/graph/badge.svg)](https://codecov.io/gh/eiko-team/eiko)
[![Docker Pulls](https://img.shields.io/docker/pulls/eikoapp/eiko?logo=docker)](https://hub.docker.com/r/eikoapp/eiko/)
[![Docker Stars](https://img.shields.io/docker/stars/eikoapp/eiko?logo=docker)](https://hub.docker.com/r/eikoapp/eiko/)
[![discord](https://img.shields.io/discord/621015347918536724?logo=discord)](https://discord.gg/NxuCWQs)
[![tipeee](https://img.shields.io/badge/tipeee-Tip!-green)](https://tipeee.com/eiko-app)
[![Godoc](https://img.shields.io/badge/godoc-eiko-blue?logo=go)](https://godoc.org/github.com/eiko-team/eiko)

eiko web app

# Installation
> Prerequirement: have golang>=1.10.4 install on your machine

You need to have a services account for the google Datastore

## Services account
```bash
export ACCOUNT_NAME=
export PROJECT_ID=
export CREDENTIALS=CREDENTIALS.json
```

```bash
gcloud iam service-accounts create $ACCOUNT_NAME
gcloud projects add-iam-policy-binding $PROJECT_ID --member "serviceAccount:$ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" --role "roles/owner"
gcloud iam service-accounts keys create $CREDENTIALS --iam-account $ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com

```

## Compilation and local building
```bash
git clone https://github.com/eiko-team/eiko.git  $GOPATH/src/github.com/eiko-team/eiko
cd $GOPATH/src/github.com/eiko-team/eiko
get get ./...
make up
```

### Display test coverage
```bash
make cover
$BROWSER test.html
```

## Docker compose
.env file content:
```
PROJECT_ID=
CREDENTIALS=CREDENTIALS.json
SALT=
```

```
  eiko:
    image: eikoapp/eiko:latest-prod
    restart: always
    environment:
      - 'PROJECT_ID=${PROJECT_ID}'
      - 'PORT=80'
      - 'STATIC_PWD=/srv'
      - 'GOOGLE_APPLICATION_CREDENTIALS=/srv/${CREDENTIALS}'
      - 'GRPC_GO_LOG_SEVERITY_LEVEL=INFO' # google datastore debug
      - 'SALT=${SALT}
    volumes:
      - './CREDENTIALS.json:/srv/CREDENTIALS.json'
      - '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt'
```

## minimified version
### Build
```bash
make mini
```

### Use in docker compose
```
  eiko:
    ...
    environment:
      - 'FILE_TYPE=minimified'
```

## Requirements
### Build minimified version
[html-minifier](https://www.npmjs.com/package/html-minifier)
[css-minifier](https://www.npmjs.com/package/uglifycss)