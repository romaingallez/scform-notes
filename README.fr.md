# SCForm Notes

[English version](README.md)

Une application web construite avec Go (Fiber) et des technologies frontend modernes pour gérer et traiter les données de formulaires. L'application utilise HTMX pour les interactions dynamiques et Tailwind CSS pour le style.

## 🚀 Fonctionnalités

- Interface web moderne avec Tailwind CSS et DaisyUI
- Rendu côté serveur avec les templates Go
- Mises à jour dynamiques de l'interface utilisateur avec HTMX
- Traitement et gestion des données de formulaires
- Système de gestion des assets
- Support de configuration d'environnement
- Rechargement en direct pendant le développement avec Air

## 🛠 Stack Technique

### Backend
- Go (framework web Fiber)
- Templates HTML
- Configuration d'environnement avec godotenv
- Air (Rechargement en direct)

### Frontend
- HTMX pour les interactions dynamiques
- Hyperscript pour une interactivité améliorée
- Tailwind CSS avec les composants DaisyUI
- Webpack pour le bundling des assets

## 📦 Prérequis

- Go 1.x
- Node.js et pnpm
- Air (outil de rechargement en direct pour Go)
- Variables d'environnement (copier depuis .env.example)

## 🚀 Mise en Route

1. Cloner le dépôt
2. Copier la configuration d'environnement :
   ```bash
   cp .env.example .env
   ```

3. Installer les dépendances frontend :
   ```bash
   pnpm install
   ```

4. Construire les assets frontend :
   ```bash
   pnpm run build
   ```

5. Exécuter l'application :
   ```bash
   # Utilisation standard de Go
   go run main.go

   # Utilisation d'Air pour le rechargement en direct pendant le développement
   air
   ```

L'application sera disponible à l'adresse `http://localhost:3000`

## 🐳 Déploiement Docker

Vous pouvez également exécuter l'application avec Docker :

1. Construire l'image Docker :
   ```bash
   docker build -t scform-notes .
   ```

2. Créer un fichier `.env.docker` à partir de l'exemple :
   ```bash
   cp .env.docker.example .env.docker
   ```
   
   Assurez-vous de mettre à jour les variables d'environnement dans `.env.docker` selon vos besoins.

3. Exécuter le conteneur browserless/chrome (requis pour le traitement des formulaires) :
   ```bash
   docker run -d -p 1337:3000 --rm --name chrome browserless/chrome
   ```

4. Exécuter le conteneur de l'application :
   ```bash
   docker run -d -p 3000:3000 --rm --env-file .env.docker --name scform-notes scform-notes
   ```

5. Accéder à l'application à l'adresse `http://localhost:3000`

### Docker Compose (Alternative)

Vous pouvez également utiliser Docker Compose pour exécuter les deux conteneurs :

1. Créer un fichier `docker-compose.yml` :
   ```yaml
   version: '3'
   services:
     app:
       build: .
       ports:
         - "3000:3000"
       env_file:
         - .env.docker
       depends_on:
         - chrome
     chrome:
       image: browserless/chrome
       ports:
         - "1337:3000"
   ```

2. Exécuter avec Docker Compose :
   ```bash
   docker-compose up -d
   ```

## 🔧 Développement

### Développement Frontend
- Surveiller les changements de Tailwind CSS :
  ```bash
  pnpm run watch
  ```

### Développement Backend
- L'application utilise les modules Go pour la gestion des dépendances
- Le point d'entrée principal de l'application se trouve dans `main.go`
- Pour le rechargement en direct pendant le développement, utilisez Air :
  ```bash
  # Installer Air si vous ne l'avez pas déjà fait
  go install github.com/cosmtrek/air@latest

  # Exécuter l'application avec Air
  air
  ```
- La logique de base est organisée dans le répertoire `internals` :
  - `scform/` : Fonctionnalités liées aux formulaires
  - `utils/` : Fonctions utilitaires et helpers
  - `web/` : Serveur web et logique de routage

## 📁 Structure du Projet

```
.
├── assets/          # Assets frontend
├── internals/       # Logique de base de l'application
│   ├── scform/      # Traitement des formulaires
│   ├── utils/       # Fonctions utilitaires
│   └── web/         # Serveur web et routage
├── views/           # Templates HTML
├── main.go         # Point d'entrée de l'application
├── go.mod          # Dépendances Go
└── package.json    # Dépendances frontend
```

## 📄 Licence

Licence ISC

## 🤝 Contribution

Les contributions, les problèmes et les demandes de fonctionnalités sont les bienvenus !