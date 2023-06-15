# React Go SSR

React Go SSR est un générateur de fichiers HTML permettant un rendu côté serveur (SSR) rapide pour les projets React.

## Introduction

React Go SSR a été développé dans le but de fournir une alternative performante à ReactSnap pour le rendu côté serveur (SSR). En générant des fichiers HTML pré-rendus, cette solution offre des performances optimales lors du chargement des pages d'un projet React.

## Utilisation en local

Pour utiliser React Go SSR, suivez les étapes suivantes :

1. Installez les dépendances requises. Go 1.20 ou supérieur est requis ainsi Chomium ou Google Chome

```bash
go mod tidy
```
2. Exécutez le programme.

```bash
go run main.go -h http://localhost:3000 -o outhtml -w 2s
```

## Utilisation en production

### React
```tsx
const rootElement = document.getElementById("root") as HTMLElement;

if (rootElement.hasChildNodes()) {
    ReactDOM.hydrateRoot(rootElement, <App />);
} else {
    ReactDOM.createRoot(rootElement).render(<App />);
}
```

### Dockerfile
```Dockerfile
# Build environment
FROM node:18 as builder
WORKDIR /usr/src/app

ENV PATH /usr/src/app/node_modules/.bin:$PATH
COPY package.* yarn.* ./
RUN npm install
COPY . ./
RUN npm run build

# SSR STUFF
# - Installer chromium
RUN apt-get update && apt-get install -y chromium
# - Téléchargement et extraction de React Go SSR
RUN curl -LO https://github.com/Squirrel-Entreprise/react-go-ssr/releases/download/v1.1.4/react-go-ssr_1.1.4_linux_amd64.tar.gz \
    && tar -xzf react-go-ssr_1.1.4_linux_amd64.tar.gz \
    && rm react-go-ssr_1.1.4_linux_amd64.tar.gz
# - Rendre le binaire exécutable
RUN chmod +x ./react-go-ssr
# - Installer serve pour rendre les fichiers React
RUN npm i -g serve
# Demarrer le serveur temporairement
RUN serve -s build -l 3000 &
# - Lancer la generation des fichiers SSR
RUN ./react-go-ssr -h http://localhost:3000 -o build -w 3s
# - Arreter le serveur
RUN pkill -f serve
# END SSR STUFF

# Production environment
FROM nginx
RUN rm -rf /etc/nginx/conf.d
RUN mkdir -p /etc/nginx/conf.d
COPY ./default.conf /etc/nginx/conf.d/
COPY --from=builder /usr/src/app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## Contribuer

Les contributions sont les bienvenues ! Si vous souhaitez contribuer à React Go SSR, veuillez envoyer vos pull requests sur la branche `master`.

## Licence

React Go SSR est distribué sous la licence MIT. Veuillez consulter le fichier `LICENSE` pour plus d'informations.