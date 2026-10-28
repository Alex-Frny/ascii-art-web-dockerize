# ASCII Art Web Dockerize

## Description

ASCII Art Web Dockerize est une application web développée en Go.

Elle permet à un utilisateur de saisir un texte, de choisir une bannière et d'afficher le résultat sous forme d'ASCII Art.

L'application est compilée dans une image Docker, puis exécutée dans un conteneur.

## Fonctionnalités

- Saisie d'un texte depuis un formulaire HTML.
- Génération d'un ASCII Art.
- Choix entre trois bannières :
  - standard
  - shadow
  - thinkertoy
- Affichage du résultat dans la page web.
- Interface responsive avec CSS.
- Gestion des méthodes HTTP.
- Gestion des erreurs HTTP 400, 404 et 500.
- Exécution de l'application dans un conteneur Docker.
- Ajout de métadonnées aux objets Docker.
- Utilisation d'un utilisateur non-root dans le conteneur.
- Gestion et nettoyage des objets Docker inutilisés.

## Technologies utilisées

- Go 1.26.5
- HTML
- CSS
- HTTP
- Docker
- Alpine Linux

Seuls les packages de la bibliothèque standard Go sont utilisés.

## Structure du projet

```text
ascii-art-web-dockerize/
├── banners/
│   ├── shadow.txt
│   ├── standard.txt
│   └── thinkertoy.txt
├── static/
│   └── style.css
├── templates/
│   └── index.html
├── .dockerignore
├── Dockerfile
├── README.md
├── go.mod
├── main.go
└── main_test.go
```

## Fonctionnement de l'application

### Route principale

```text
GET /
```

Cette route affiche la page principale et le formulaire permettant de saisir un texte et de choisir une bannière.

### Route de génération

```text
POST /ascii-art
```

Cette route récupère les données du formulaire et génère l'ASCII Art.

Les données reçues sont :

- `text` : le texte à transformer ;
- `banner` : la bannière à utiliser.

## Statuts HTTP

| Situation | Statut |
|---|---:|
| Page principale correctement affichée | 200 |
| ASCII Art correctement généré | 200 |
| Texte vide | 400 |
| Bannière invalide | 400 |
| Mauvaise méthode HTTP | 400 |
| Route inexistante | 404 |
| Erreur interne ou fichier manquant | 500 |

## Lancement sans Docker

Depuis la racine du projet, lancer :

```bash
go run .
```

Ouvrir ensuite dans un navigateur :

```text
http://localhost:8080
```

## Tests Go

Le fichier `main_test.go` contient les tests du projet.

Pour lancer tous les tests :

```bash
go test -v
```

Les tests vérifient notamment :

- la route principale ;
- les routes inexistantes ;
- les méthodes HTTP incorrectes ;
- les textes vides ;
- les bannières invalides ;
- les caractères non pris en charge ;
- la génération de l'ASCII Art ;
- la gestion des retours à la ligne ;
- la lecture des fichiers de bannière.

Le résultat attendu se termine par :

```text
PASS
ok
```

## Dockerfile

Le Dockerfile utilise une construction en deux étapes appelée `multi-stage build`.

### Première étape : compilation

L'image Go est utilisée pour compiler le programme :

```dockerfile
FROM golang:1.26.5-alpine AS builder
```

Le programme est compilé sous la forme d'un exécutable Linux nommé `server`.

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o server .
```

### Deuxième étape : exécution

Une image Alpine plus légère est utilisée pour exécuter uniquement le serveur compilé :

```dockerfile
FROM alpine:3.22
```

Le compilateur Go n'est donc pas conservé dans l'image finale. Cela permet de réduire la taille de l'image.

Un utilisateur non-root nommé `appuser` est également créé afin de ne pas exécuter le serveur avec les droits administrateur.

Le port utilisé par le serveur est le port `8080`.

```dockerfile
EXPOSE 8080
```

La commande suivante démarre le serveur lorsque le conteneur est lancé :

```dockerfile
CMD ["./server"]
```

## Construction de l'image Docker

Depuis la racine du projet, construire l'image avec :

```bash
docker build -t ascii-art-web:1.0 .
```

Signification de la commande :

- `docker build` construit une image ;
- `-t` attribue un nom et une version à l'image ;
- `ascii-art-web` est le nom de l'image ;
- `1.0` est la version de l'image ;
- `.` indique que le Dockerfile se trouve dans le dossier actuel.

Vérifier la présence de l'image :

```bash
docker images
```

## Métadonnées de l'image

Le Dockerfile ajoute plusieurs métadonnées à l'image :

```dockerfile
LABEL maintainer="Alexis Fourny"
LABEL version="1.0"
LABEL description="Serveur web ASCII Art développé en Go"
```

Afficher les métadonnées de l'image :

```bash
docker image inspect ascii-art-web:1.0 --format '{{json .Config.Labels}}'
```

## Création du conteneur

Créer et lancer le conteneur avec :

```bash
docker run -d \
  --name ascii-art-web-container \
  --label project=ascii-art-web-dockerize \
  --label author="Alexis Fourny" \
  -p 8080:8080 \
  ascii-art-web:1.0
```

Signification de la commande :

- `-d` exécute le conteneur en arrière-plan ;
- `--name` attribue un nom au conteneur ;
- `--label` ajoute des métadonnées ;
- `-p 8080:8080` relie le port local au port du conteneur ;
- `ascii-art-web:1.0` indique l'image à utiliser.

## Vérification du conteneur

Afficher les conteneurs actifs :

```bash
docker ps
```

Le conteneur doit avoir le statut `Up`.

Afficher tous les conteneurs, y compris les conteneurs arrêtés :

```bash
docker ps -a
```

Afficher les métadonnées du conteneur :

```bash
docker container inspect ascii-art-web-container --format '{{json .Config.Labels}}'
```

Afficher les journaux du serveur :

```bash
docker logs ascii-art-web-container
```

## Utilisation avec Docker

Lorsque le conteneur est actif, ouvrir dans un navigateur :

```text
http://localhost:8080
```

L'utilisateur peut ensuite :

1. saisir un texte ;
2. choisir une bannière ;
3. envoyer le formulaire ;
4. consulter le résultat en ASCII Art.

## Tests HTTP

### Tester la page principale

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/
```

Résultat attendu :

```text
200
```

### Tester une route inexistante

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/inconnue
```

Résultat attendu :

```text
404
```

### Tester une génération valide

```bash
curl -s -o /dev/null -w "%{http_code}\n" \
  -X POST \
  -d "text=Hello&banner=standard" \
  http://localhost:8080/ascii-art
```

Résultat attendu :

```text
200
```

### Tester un texte vide

```bash
curl -s -o /dev/null -w "%{http_code}\n" \
  -X POST \
  -d "text=&banner=standard" \
  http://localhost:8080/ascii-art
```

Résultat attendu :

```text
400
```

### Tester une bannière invalide

```bash
curl -s -o /dev/null -w "%{http_code}\n" \
  -X POST \
  -d "text=Hello&banner=inconnue" \
  http://localhost:8080/ascii-art
```

Résultat attendu :

```text
400
```

## Gestion du conteneur

Arrêter le conteneur :

```bash
docker stop ascii-art-web-container
```

Redémarrer le conteneur :

```bash
docker start ascii-art-web-container
```

Supprimer le conteneur après l'avoir arrêté :

```bash
docker rm ascii-art-web-container
```

Il est également possible de l'arrêter et de le supprimer avec une seule commande :

```bash
docker rm -f ascii-art-web-container
```

## Gestion des objets Docker inutilisés

Docker peut conserver des images, des conteneurs et des caches devenus inutiles.

Afficher l'espace utilisé par Docker :

```bash
docker system df
```

Afficher les conteneurs arrêtés :

```bash
docker ps -a
```

Supprimer les conteneurs arrêtés inutilisés :

```bash
docker container prune
```

Supprimer les images inutilisées :

```bash
docker image prune
```

Supprimer le cache de construction inutilisé :

```bash
docker builder prune
```

La commande suivante doit être utilisée avec prudence :

```bash
docker system prune -a
```

Elle peut supprimer :

- les conteneurs arrêtés ;
- les images non utilisées ;
- les réseaux inutilisés ;
- le cache de construction.

Il est recommandé d'utiliser `docker system df` et `docker ps -a` avant de lancer une commande de nettoyage.

## Arrêt complet du projet

Arrêter le conteneur :

```bash
docker stop ascii-art-web-container
```

Supprimer le conteneur :

```bash
docker rm ascii-art-web-container
```

Supprimer l'image du projet :

```bash
docker rmi ascii-art-web:1.0
```

## Bonnes pratiques appliquées

- Utilisation d'un Dockerfile.
- Utilisation d'une construction en plusieurs étapes.
- Utilisation d'une image Alpine légère.
- Compilation du programme avec `CGO_ENABLED=0`.
- Utilisation d'un utilisateur non-root.
- Copie uniquement des fichiers nécessaires dans l'image finale.
- Déclaration du port avec `EXPOSE`.
- Ajout de métadonnées avec `LABEL`.
- Utilisation d'un fichier `.dockerignore`.
- Attribution d'un nom et d'une version à l'image.
- Attribution d'un nom et de labels au conteneur.
- Gestion des objets Docker inutilisés.

## Auteur

Alexis Fourny





# Guide de passage d'audit - ASCII Art Web Dockerize

Ce document contient les commandes à exécuter pendant l'audit du projet.

 docker image build -f Dockerfile -t ascii-art-web:1.0 .

  docker images

  docker container run -p 8080:8080 --detach --name ascii-art-web-container ascii-art-web:1.0

   docker ps -a --filter name=ascii-art-web-container

   docker exec -it ascii-art-web-container /bin/sh

   docker image inspect   ascii-art-web:1.0   --format '{{json .Config.Labels}}'

   ./build.sh

   