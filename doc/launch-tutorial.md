# Comment faire fonctionner Eiko sur sa machine
<walkthrough-author name="Tom Moulard"></walkthrough-author>
<walkthrough-tutorial-duration duration=2></walkthrough-tutorial-duration>
<walkthrough-devshell-precreate></walkthrough-devshell-precreate>

## Introduction
Ce guide va vous montrer comment démarrer le projet en l'installant puis en démarrant le serveur.

Il faut avoir au minimum [Golang](https://golang.org/doc/install) avec une version >= à 1.13.4 d'installé sur la machine.

## Installation
Pour installer le projet sur votre machine, il faut faire:
```bash
go get github.com/eiko-team/eiko
```
ou
```bash
git clone https://github.com/eiko-team/eiko.git
```
(sur Cloud shell cette  étape n'est pas nécessaire)

### Mise en place de l'environnement
Pour que le projet fonctionne, il faut stocker des informations dans l'environnement pour la configuration de l'application:

```bash
export PORT=8080
export SALT="{{project-name}}"
export STATIC_PWD=$(pwd)
export SEARCH_APP_ID=""
export SEARCH_API_KEY=""
export PROJECT_ID="{{project-id}}"
export GOOGLE_APPLICATION_CREDENTIALS=""
```

Où:
 - `PORT`: Le port sur lequel l'appication va se connecter.
 - `SALT`: Le sel utilisé pour hasher les mot de passe.
 - `STATIC_PWD`: le dossier racine des fichiers source html.
 - `SEARCH_APP_ID`: Identifiant de l'application d'Algolia.
 - `SEARCH_API_KEY`: Clef d'API d'Algolia.
 - `PROJECT_ID`: the GCP(Google Cloud Platform) Project id
 - `GOOGLE_APPLICATION_CREDENTIALS`: Chemin d'accès au fichier de configuration de GCP.

<walkthrough-editor-open-file filePath="README.md"
                                text="Voir comment créer le fichier de configuration.">
</walkthrough-editor-open-file>

## Compilation
Pour voir si le projet peux fonctionner sur votre machine, il faut faire
```bash
make
```

### Lancer le serveur
Via Docker:
```bash
make up
```

Via un shell:
```bash
./app
```

## Profiter
Voila le serveur lancé, [il n'y a plus qu'à utiliser une navigateur pour accéder à la page d’accueil](http://127.0.0.1).

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>