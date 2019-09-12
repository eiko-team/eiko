# eiko
[![Go Report Card](https://goreportcard.com/badge/github.com/eiko-team/eiko)](https://goreportcard.com/report/github.com/eiko-team/eiko)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/13cbb61d7e734f16a8f0494e0a13a993)](https://www.codacy.com/manual/tomMoulard/eiko?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=eiko-team/eiko&amp;utm_campaign=Badge_Grade)

eiko web app

## Compte de services

```bash
export NAME= #account name
export PROJECT_ID=
export CREDENTIALS=CREDENTIALS.json
```

```bash
$ gcloud iam service-accounts create $NAME
$ gcloud projects add-iam-policy-binding $PROJECT_ID --member "serviceAccount:$NAME@$PROJECT_ID.iam.gserviceaccount.com" --role "roles/owner"
$ gcloud iam service-accounts keys create $CREDENTIALS --iam-account $NAME@$PROJECT_ID.iam.gserviceaccount.com
