# FactuGest-WebInformatique - Documentation Complète

## Table des matières

1. [Présentation du projet](#présentation-du-projet)
2. [Architecture technique](#architecture-technique)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Utilisation](#utilisation)
6. [Fonctionnalités](#fonctionnalités)
7. [API Documentation](#api-documentation)
8. [Maintenance](#maintenance)
9. [Dépannage](#dépannage)
10. [Support](#support)

---

## Présentation du projet

**FactuGest-WebInformatique** est un logiciel web complet de gestion de factures et de clients développé spécifiquement pour l'entreprise WebInformatique-Sarl. Ce système permet de gérer efficacement les clients, créer et suivre les factures, et obtenir des statistiques en temps réel sur l'activité commerciale.

### Caractéristiques principales

- **Interface web moderne et responsive** : Accessible depuis tout navigateur web
- **Gestion complète des clients** : Création, modification, suppression avec toutes les informations importantes
- **Gestion des factures** : Création, suivi, modification avec calculs automatiques de TVA
- **Tableau de bord administrateur** : Statistiques en temps réel et activité récente
- **Synchronisation multi-utilisateurs** : Accès simultané depuis plusieurs postes
- **Design professionnel** : Interface intuitive avec animations et effets visuels
- **Sécurité** : Authentification et protection des données

### Technologies utilisées

- **Backend** : Go (Golang) avec framework Gorilla Mux
- **Base de données** : MySQL 8.0
- **Frontend** : HTML5, CSS3, JavaScript (Vanilla)
- **Serveur web** : Serveur HTTP intégré Go
- **Authentification** : JWT (JSON Web Tokens)

---

## Architecture technique

### Structure du projet

```
FactuGest-WebInformatique/
├── backend/                 # Code source du serveur Go
│   ├── main.go             # Point d'entrée de l'application
│   ├── config/             # Configuration de l'application
│   ├── database/           # Connexion et gestion de la base de données
│   ├── handlers/           # Gestionnaires des routes API
│   ├── models/             # Modèles de données
│   └── go.mod              # Dépendances Go
├── frontend/               # Interface utilisateur web
│   ├── index.html          # Page principale
│   ├── styles.css          # Feuilles de style
│   └── script.js           # Logique JavaScript
├── database/               # Scripts de base de données
│   └── init.sql            # Script d'initialisation
└── docs/                   # Documentation
    └── README.md           # Ce fichier
```

### Base de données

La base de données MySQL contient les tables suivantes :

- **clients** : Informations des clients (nom, prénom, adresse, contact, etc.)
- **factures** : Données des factures (numéro, montants, dates, statut)
- **lignes_factures** : Détails des lignes de factures (description, quantité, prix)
- **utilisateurs** : Comptes utilisateurs pour l'authentification

### API REST

L'API suit les conventions REST avec les endpoints suivants :

- `GET /api/clients` : Liste des clients
- `POST /api/clients` : Création d'un client
- `PUT /api/clients/{id}` : Modification d'un client
- `DELETE /api/clients/{id}` : Suppression d'un client
- `GET /api/factures` : Liste des factures
- `POST /api/factures` : Création d'une facture
- `PUT /api/factures/{id}` : Modification d'une facture
- `DELETE /api/factures/{id}` : Suppression d'une facture
- `GET /api/admin/stats` : Statistiques administrateur

---


## Installation

### Prérequis système

- **Système d'exploitation** : Linux (Ubuntu 20.04+ recommandé), Windows 10+, ou macOS 10.15+
- **Go** : Version 1.21 ou supérieure
- **MySQL** : Version 8.0 ou supérieure
- **Navigateur web** : Chrome, Firefox, Safari, ou Edge (versions récentes)

### Étapes d'installation

#### 1. Installation de Go

**Sur Ubuntu/Debian :**
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Sur Windows :**
Télécharger l'installateur depuis https://golang.org/dl/ et suivre les instructions.

#### 2. Installation de MySQL

**Sur Ubuntu/Debian :**
```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql
```

**Sur Windows :**
Télécharger MySQL Community Server depuis https://dev.mysql.com/downloads/

#### 3. Configuration de la base de données

```bash
# Se connecter à MySQL en tant que root
sudo mysql

# Exécuter le script d'initialisation
mysql> source /chemin/vers/FactuGest-WebInformatique/database/init.sql;
mysql> exit;
```

#### 4. Installation des dépendances Go

```bash
cd FactuGest-WebInformatique/backend
go mod tidy
```

#### 5. Compilation du serveur

```bash
go build -o factugest-server .
```

### Configuration

#### Variables d'environnement

Le système utilise les variables d'environnement suivantes (avec valeurs par défaut) :

```bash
export PORT=8080                    # Port du serveur web
export DB_HOST=localhost            # Adresse du serveur MySQL
export DB_PORT=3306                 # Port MySQL
export DB_USER=factugest           # Utilisateur MySQL
export DB_PASS=factugest123        # Mot de passe MySQL
export DB_NAME=factugest_db        # Nom de la base de données
export JWT_SECRET=factugest-secret-key-2024  # Clé secrète JWT
```

#### Fichier de configuration

Vous pouvez également modifier les paramètres dans le fichier `backend/config/config.go` si nécessaire.

### Démarrage du système

#### Démarrage manuel

```bash
cd FactuGest-WebInformatique/backend
./factugest-server
```

Le serveur sera accessible à l'adresse : http://localhost:8080

#### Démarrage automatique (Linux)

Créer un service systemd :

```bash
sudo nano /etc/systemd/system/factugest.service
```

Contenu du fichier :
```ini
[Unit]
Description=FactuGest WebInformatique Server
After=network.target mysql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/chemin/vers/FactuGest-WebInformatique/backend
ExecStart=/chemin/vers/FactuGest-WebInformatique/backend/factugest-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Activer et démarrer le service :
```bash
sudo systemctl daemon-reload
sudo systemctl enable factugest
sudo systemctl start factugest
```

---


## Utilisation

### Accès à l'application

1. Ouvrir un navigateur web
2. Aller à l'adresse : http://localhost:8080 (ou l'adresse IP du serveur)
3. L'interface de FactuGest s'affiche automatiquement

### Connexion

Par défaut, un compte administrateur est créé :
- **Nom d'utilisateur** : admin
- **Mot de passe** : password

### Navigation dans l'interface

L'interface est organisée en 4 onglets principaux :

#### 1. Toutes les factures
- Affichage de toutes les factures créées
- Filtrage par statut (brouillon, envoyée, payée, en retard)
- Recherche par numéro de facture ou nom de client
- Actions : modifier, supprimer, générer PDF

#### 2. Gestion des clients
- Liste de tous les clients enregistrés
- Recherche par nom, prénom ou entreprise
- Actions : modifier, supprimer

#### 3. Créer un client
- Formulaire de création d'un nouveau client
- Champs obligatoires : nom et prénom
- Champs optionnels : entreprise, SIRET, adresse, contact

#### 4. Administration
- Statistiques générales (nombre de clients, factures, chiffre d'affaires)
- Activité récente
- Gestion des utilisateurs (pour les administrateurs)

### Workflow typique

#### Création d'un client
1. Aller dans l'onglet "Créer un client"
2. Remplir les informations obligatoires (nom, prénom)
3. Ajouter les informations complémentaires si nécessaire
4. Cliquer sur "Enregistrer le client"

#### Création d'une facture
1. Aller dans l'onglet "Toutes les factures"
2. Cliquer sur "Nouvelle facture"
3. Sélectionner un client dans la liste déroulante
4. Remplir les informations de la facture :
   - Numéro de facture (généré automatiquement)
   - Date de facture et d'échéance
   - Montant HT
   - Taux de TVA (20% par défaut)
   - Description et notes
5. Cliquer sur "Enregistrer"

#### Suivi des factures
1. Dans l'onglet "Toutes les factures", voir l'état de chaque facture
2. Utiliser les filtres pour afficher seulement certains statuts
3. Modifier le statut d'une facture en la modifiant

---

## Fonctionnalités

### Gestion des clients

#### Informations stockées
- **Informations personnelles** : Nom, prénom
- **Informations professionnelles** : Entreprise, SIRET
- **Coordonnées** : Adresse complète, téléphone, email
- **Métadonnées** : Date de création, dernière modification

#### Fonctionnalités
- **Création** : Formulaire complet avec validation
- **Modification** : Édition de toutes les informations
- **Suppression** : Avec vérification des factures associées
- **Recherche** : Par nom, prénom ou entreprise
- **Export** : Données au format JSON via API

### Gestion des factures

#### Informations stockées
- **Identification** : Numéro unique, client associé
- **Dates** : Date de facture, date d'échéance
- **Montants** : HT, TVA, TTC (calculs automatiques)
- **Statut** : Brouillon, envoyée, payée, en retard
- **Contenu** : Description, notes

#### Fonctionnalités
- **Création** : Avec numérotation automatique
- **Modification** : Tous les champs modifiables
- **Suppression** : Avec confirmation
- **Filtrage** : Par statut, client, période
- **Recherche** : Par numéro ou nom de client
- **Calculs automatiques** : TVA et montant TTC

### Tableau de bord administrateur

#### Statistiques affichées
- **Nombre total de clients** : Compteur en temps réel
- **Nombre total de factures** : Toutes factures confondues
- **Chiffre d'affaires** : Somme des montants TTC
- **Répartition par statut** : Graphique des statuts de factures

#### Activité récente
- **5 dernières factures** : Avec client et montant
- **Mise à jour en temps réel** : Actualisation automatique

### Sécurité et authentification

#### Authentification
- **JWT (JSON Web Tokens)** : Tokens sécurisés avec expiration
- **Sessions persistantes** : Reconnexion automatique
- **Déconnexion sécurisée** : Invalidation des tokens

#### Sécurité des données
- **Validation côté serveur** : Toutes les entrées sont validées
- **Protection CORS** : Configuration des origines autorisées
- **Chiffrement des mots de passe** : Utilisation de bcrypt
- **Protection SQL** : Requêtes préparées contre l'injection SQL

### Interface utilisateur

#### Design et ergonomie
- **Design moderne** : Interface avec dégradés et animations
- **Responsive** : Adaptation automatique aux écrans mobiles
- **Accessibilité** : Navigation au clavier et lecteurs d'écran
- **Feedback visuel** : Notifications et états de chargement

#### Fonctionnalités UX
- **Navigation par onglets** : Interface intuitive
- **Filtres en temps réel** : Recherche instantanée
- **Modales** : Formulaires en overlay
- **Animations** : Transitions fluides entre les états

---


## API Documentation

### Authentification

Toutes les requêtes API (sauf login) nécessitent un token JWT dans l'en-tête :
```
Authorization: Bearer <token>
```

### Endpoints disponibles

#### Authentification
```http
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

**Réponse :**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    }
  }
}
```

#### Clients

**Lister tous les clients :**
```http
GET /api/clients
```

**Créer un client :**
```http
POST /api/clients
Content-Type: application/json

{
  "nom": "Dupont",
  "prenom": "Jean",
  "entreprise": "Dupont SARL",
  "email": "jean.dupont@email.com",
  "telephone": "01.23.45.67.89",
  "adresse": "123 Rue de la Paix",
  "ville": "Paris",
  "code_postal": "75001"
}
```

**Modifier un client :**
```http
PUT /api/clients/{id}
Content-Type: application/json

{
  "nom": "Dupont",
  "prenom": "Jean-Pierre",
  "entreprise": "Dupont & Fils SARL"
}
```

**Supprimer un client :**
```http
DELETE /api/clients/{id}
```

#### Factures

**Lister toutes les factures :**
```http
GET /api/factures
```

**Créer une facture :**
```http
POST /api/factures
Content-Type: application/json

{
  "numero_facture": "FACT-2024-0001",
  "client_id": 1,
  "date_facture": "2024-01-15",
  "date_echeance": "2024-02-15",
  "montant_ht": 1000.00,
  "taux_tva": 20.00,
  "description": "Prestation de développement web",
  "statut": "brouillon"
}
```

#### Administration

**Statistiques :**
```http
GET /api/admin/stats
```

**Réponse :**
```json
{
  "success": true,
  "data": {
    "total_clients": 15,
    "total_factures": 42,
    "montant_total": 25000.00,
    "factures_par_statut": {
      "brouillon": 5,
      "envoyee": 20,
      "payee": 15,
      "en_retard": 2
    },
    "factures_recentes": [...]
  }
}
```

### Codes de réponse

- **200 OK** : Requête réussie
- **201 Created** : Ressource créée avec succès
- **400 Bad Request** : Données invalides
- **401 Unauthorized** : Authentification requise
- **404 Not Found** : Ressource non trouvée
- **409 Conflict** : Conflit (ex: suppression impossible)
- **500 Internal Server Error** : Erreur serveur

---

## Maintenance

### Sauvegarde de la base de données

#### Sauvegarde complète
```bash
mysqldump -u factugest -p factugest_db > backup_$(date +%Y%m%d_%H%M%S).sql
```

#### Sauvegarde automatique (crontab)
```bash
# Ajouter dans crontab (crontab -e)
0 2 * * * mysqldump -u factugest -p factugest_db > /backups/factugest_$(date +\%Y\%m\%d).sql
```

#### Restauration
```bash
mysql -u factugest -p factugest_db < backup_20240115_020000.sql
```

### Mise à jour du logiciel

1. **Arrêter le service :**
```bash
sudo systemctl stop factugest
```

2. **Sauvegarder la base de données :**
```bash
mysqldump -u factugest -p factugest_db > backup_avant_maj.sql
```

3. **Mettre à jour le code :**
```bash
cd FactuGest-WebInformatique
git pull origin main  # Si utilisation de Git
```

4. **Recompiler :**
```bash
cd backend
go build -o factugest-server .
```

5. **Redémarrer le service :**
```bash
sudo systemctl start factugest
```

### Monitoring

#### Vérification du statut
```bash
sudo systemctl status factugest
```

#### Logs du service
```bash
sudo journalctl -u factugest -f
```

#### Vérification de la base de données
```bash
mysql -u factugest -p -e "SELECT COUNT(*) FROM factugest_db.clients;"
```

---

## Dépannage

### Problèmes courants

#### Le serveur ne démarre pas

**Symptôme :** Erreur au démarrage du service

**Solutions :**
1. Vérifier que MySQL est démarré :
```bash
sudo systemctl status mysql
```

2. Vérifier les paramètres de connexion dans `config/config.go`

3. Tester la connexion à la base :
```bash
mysql -u factugest -p factugest_db
```

#### Erreur "Connexion refusée"

**Symptôme :** Impossible d'accéder à l'interface web

**Solutions :**
1. Vérifier que le serveur écoute sur le bon port :
```bash
netstat -tlnp | grep 8080
```

2. Vérifier le firewall :
```bash
sudo ufw status
sudo ufw allow 8080
```

#### Erreurs de base de données

**Symptôme :** Erreurs lors des opérations CRUD

**Solutions :**
1. Vérifier les logs :
```bash
sudo journalctl -u factugest -n 50
```

2. Vérifier l'espace disque :
```bash
df -h
```

3. Redémarrer MySQL :
```bash
sudo systemctl restart mysql
```

#### Interface lente ou non responsive

**Solutions :**
1. Vider le cache du navigateur
2. Vérifier la charge du serveur :
```bash
top
htop
```

3. Optimiser MySQL :
```sql
OPTIMIZE TABLE clients, factures;
```

### Logs et diagnostics

#### Localisation des logs
- **Logs système :** `/var/log/syslog`
- **Logs MySQL :** `/var/log/mysql/error.log`
- **Logs application :** `journalctl -u factugest`

#### Activation du mode debug
Modifier `main.go` pour ajouter plus de logs :
```go
log.SetLevel(log.DebugLevel)
```

---

## Support

### Informations système

Pour toute demande de support, fournir les informations suivantes :

```bash
# Version du système
cat /etc/os-release

# Version de Go
go version

# Version de MySQL
mysql --version

# Statut des services
sudo systemctl status factugest mysql

# Logs récents
sudo journalctl -u factugest -n 20
```

### Contact

- **Développeur :** Assistant IA Manus
- **Entreprise :** WebInformatique-Sarl
- **Documentation :** Ce fichier README.md

### Ressources utiles

- **Documentation Go :** https://golang.org/doc/
- **Documentation MySQL :** https://dev.mysql.com/doc/
- **Gorilla Mux :** https://github.com/gorilla/mux

---

## Licence et crédits

**FactuGest-WebInformatique** - Logiciel de gestion de factures et clients

Développé pour WebInformatique-Sarl par l'assistant IA Manus.

Ce logiciel est fourni "tel quel" sans garantie d'aucune sorte. L'utilisation se fait aux risques et périls de l'utilisateur.

---

*Documentation générée le $(date)*
*Version du logiciel : 1.0.0*

