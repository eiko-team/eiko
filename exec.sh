#!/bin/bash

# copy that in powershell
$env:NAME="api-eiko-test"
$env:PROJECT_ID="poetic-hexagon-252009"
$env:SALT="salt"
$env:PORT=80
$env:STATIC_PWD="$(pwd)"
$env:SEARCH_APP_ID="UVF4TXM9ZI"
$env:SEARCH_API_KEY="ec3dd566bafa202b6c96cea7f539ff20"
$env:GOOGLE_APPLICATION_CREDENTIALS="kamardine-chikh_credentials.json"

./app.exe
