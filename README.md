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
[... React build]

# SSR STUFF
# Copier le script dans le conteneur
COPY ./ssr-stuff.sh /usr/src/app/ssr-stuff.sh
# Rendre le script exécutable
RUN chmod +x /usr/src/app/ssr-stuff.sh
# Exécuter le script
RUN /usr/src/app/ssr-stuff.sh
# END SSR STUFF

[... Nginx]
```

## Contribuer

Les contributions sont les bienvenues ! Si vous souhaitez contribuer à React Go SSR, veuillez envoyer vos pull requests sur la branche `master`.

## Licence

React Go SSR est distribué sous la licence MIT. Veuillez consulter le fichier `LICENSE` pour plus d'informations.