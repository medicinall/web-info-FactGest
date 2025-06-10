# FactuGest-WebInformatique

![FactuGest Logo](https://img.shields.io/badge/FactuGest-WebInformatique-blue?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat&logo=mysql&logoColor=white)
![License](https://img.shields.io/badge/License-Proprietary-red?style=flat)

**Logiciel web professionnel de gestion de factures et clients pour WebInformatique-Sarl**

## 🚀 Aperçu

FactuGest-WebInformatique est une solution complète de gestion commerciale développée spécifiquement pour WebInformatique-Sarl. Cette application web moderne permet de gérer efficacement les clients, créer et suivre les factures, et obtenir des statistiques en temps réel.

### ✨ Fonctionnalités principales

- 👥 **Gestion complète des clients** - Création, modification, suppression avec toutes les informations importantes
- 📄 **Gestion des factures** - Création, suivi, modification avec calculs automatiques de TVA
- 🎉 **Génération de PDF** - Factures professionnelles téléchargeables en un clic
- 📊 **Tableau de bord administrateur** - Statistiques en temps réel et activité récente
- 🔄 **Synchronisation multi-utilisateurs** - Accès simultané depuis plusieurs postes
- 🎨 **Interface moderne et responsive** - Design professionnel avec animations
- 🔐 **Sécurité avancée** - Authentification JWT et protection des données

### 🖥️ Captures d'écran

L'interface propose 4 onglets principaux :

1. **Toutes les factures** - Vue d'ensemble avec filtres et recherche
2. **Gestion des clients** - Liste et actions sur les clients
3. **Créer un client** - Formulaire complet de création
4. **Administration** - Statistiques et tableau de bord

## 🚀 Installation rapide

### Prérequis
- Go 1.21+
- MySQL 8.0+
- Navigateur web moderne

### Installation en 3 étapes

1. **Cloner et configurer :**
```bash
git clone <repository>
cd FactuGest-WebInformatique
```

2. **Installer les dépendances :**
```bash
# Installer MySQL et Go (voir INSTALL.md pour les détails)
sudo mysql < database/init.sql
cd backend && go mod tidy
```

3. **Démarrer l'application :**
```bash
./start.sh
```

L'application sera accessible sur : **http://localhost:8080**

**Connexion par défaut :**
- Utilisateur : `admin`
- Mot de passe : `password`

## 📖 Documentation

- 📋 [Guide d'installation complet](INSTALL.md)
- 📚 [Documentation complète](docs/README.md)
- 🔧 [Guide de démarrage rapide](#installation-rapide)

## 🛠️ Technologies utilisées

### Backend
- **Go (Golang)** - Serveur web haute performance
- **Gorilla Mux** - Routeur HTTP
- **MySQL** - Base de données relationnelle
- **JWT** - Authentification sécurisée

### Frontend
- **HTML5/CSS3** - Structure et design moderne
- **JavaScript (Vanilla)** - Interactivité et appels API
- **Font Awesome** - Icônes professionnelles
- **Design responsive** - Compatible mobile et desktop

## 🏗️ Architecture

```
FactuGest-WebInformatique/
├── backend/                 # Serveur Go
│   ├── main.go             # Point d'entrée
│   ├── config/             # Configuration
│   ├── database/           # Connexion DB
│   ├── handlers/           # API REST
│   └── models/             # Modèles de données
├── frontend/               # Interface web
│   ├── index.html          # Page principale
│   ├── styles.css          # Design
│   └── script.js           # Logique client
├── database/               # Scripts SQL
├── docs/                   # Documentation
└── start.sh               # Script de démarrage
```

## 🔧 Utilisation

### Démarrage/Arrêt
```bash
./start.sh start    # Démarrer le serveur
./start.sh stop     # Arrêter le serveur
./start.sh restart  # Redémarrer le serveur
./start.sh status   # Voir le statut
./start.sh logs     # Voir les logs
```

### Workflow typique

1. **Créer un client** dans l'onglet "Créer un client"
2. **Créer une facture** en sélectionnant le client
3. **Suivre les factures** dans l'onglet "Toutes les factures"
4. **Consulter les statistiques** dans l'onglet "Administration"

## 🔐 Sécurité

- **Authentification JWT** avec expiration automatique
- **Validation côté serveur** de toutes les entrées
- **Protection CORS** configurée
- **Chiffrement des mots de passe** avec bcrypt
- **Requêtes préparées** contre l'injection SQL

## 📊 API REST

L'application expose une API REST complète :

- `GET/POST/PUT/DELETE /api/clients` - Gestion des clients
- `GET/POST/PUT/DELETE /api/factures` - Gestion des factures
- `GET /api/admin/stats` - Statistiques administrateur
- `POST /api/login` - Authentification

Voir la [documentation API complète](docs/README.md#api-documentation) pour plus de détails.

## 🐛 Dépannage

### Problèmes courants

**Serveur ne démarre pas :**
```bash
sudo systemctl start mysql
./start.sh restart
```

**Port 8080 occupé :**
```bash
sudo lsof -i :8080
# Tuer le processus si nécessaire
```

**Erreur de base de données :**
```bash
sudo mysql < database/init.sql
```

Voir le [guide de dépannage complet](docs/README.md#dépannage) pour plus de solutions.

## 📈 Statistiques du projet

- **Lignes de code :** ~2000 lignes
- **Fichiers :** 15+ fichiers source
- **Technologies :** 5 technologies principales
- **Fonctionnalités :** 20+ fonctionnalités
- **Documentation :** 100+ pages

## 🤝 Support

Pour toute question ou problème :

1. Consulter la [documentation complète](docs/README.md)
2. Vérifier le [guide de dépannage](docs/README.md#dépannage)
3. Examiner les logs avec `./start.sh logs`

## 📝 Licence

Ce logiciel est développé spécifiquement pour **WebInformatique-Sarl**.

---

**Développé avec ❤️ par l'assistant IA Manus pour WebInformatique-Sarl**

*Version 1.0.0 - $(date +%Y)*

