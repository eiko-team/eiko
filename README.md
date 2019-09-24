# eiko
![maintenance](https://img.shields.io/maintenance/yes/2019)
[![CircleCI](https://circleci.com/gh/eiko-team/eiko.svg?style=svg)](https://circleci.com/gh/eiko-team/eiko)
[![Go Report Card](https://goreportcard.com/badge/github.com/eiko-team/eiko)](https://goreportcard.com/report/github.com/eiko-team/eiko)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/13cbb61d7e734f16a8f0494e0a13a993)](https://www.codacy.com/manual/tomMoulard/eiko?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=eiko-team/eiko&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/eiko-team/eiko/branch/master/graph/badge.svg)](https://codecov.io/gh/eiko-team/eiko)

eiko web app

## Compte de services

```bash
export ACCOUNT_NAME=
export PROJECT_ID=
export CREDENTIALS=CREDENTIALS.json
export SALT= # For the password hashing
```

```bash
gcloud iam service-accounts create $ACCOUNT_NAME
gcloud projects add-iam-policy-binding $PROJECT_ID --member "serviceAccount:$ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" --role "roles/owner"
gcloud iam service-accounts keys create $CREDENTIALS --iam-account $ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com
