# eiko
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
