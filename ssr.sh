#!/bin/bash

# SSR STUFF
# Installer chromium
apt-get update && \
apt-get install -y \
chromium \
# chromium dependencies
libnss3 \
libxss1 \
libasound2 \
libxtst6 \
libgtk-3-0 \
libgbm1 \
ca-certificates \
# fonts
fonts-liberation fonts-noto-color-emoji fonts-noto-cjk \
# timezone
tzdata \
# process reaper
dumb-init \
# headful mode support, for example: $ xvfb-run chromium-browser --remote-debugging-port=9222
xvfb && \
# cleanup
rm -rf /var/lib/apt/lists/*
# - Téléchargement et extraction de React Go SSR
curl -LO https://github.com/Squirrel-Entreprise/react-go-ssr/releases/download/v1.1.9/react-go-ssr_1.1.9_linux_amd64.tar.gz && \
tar -xzf react-go-ssr_1.1.9_linux_amd64.tar.gz && \
rm react-go-ssr_1.1.9_linux_amd64.tar.gz
# - Rendre le binaire exécutable
chmod +x ./react-go-ssr
# - Installer serve pour rendre les fichiers React
npm i -g serve
# Demarrer le serveur temporairement
serve -s /usr/src/app/build -l 1234 &
serve_pid=$!
# Attendre que le serveur démarre
sleep 5
# - Lancer la generation des fichiers SSR
./react-go-ssr -h http://localhost:1234 -o /usr/src/app/build -w 3s
# - Arreter le serveur
kill $serve_pid
# END SSR STUFF
