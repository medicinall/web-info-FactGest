# Guide d'installation rapide - FactuGest-WebInformatique

## Installation en 5 minutes

### 1. Prérequis
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server wget

# Démarrer MySQL
sudo systemctl start mysql
sudo systemctl enable mysql
```

### 2. Installation de Go
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 3. Configuration de la base de données
```bash
sudo mysql < database/init.sql
```

### 4. Compilation et démarrage
```bash
cd backend
go mod tidy
go build -o factugest-server .
./factugest-server
```

### 5. Accès à l'application
Ouvrir un navigateur et aller à : http://localhost:8080

**Connexion par défaut :**
- Utilisateur : admin
- Mot de passe : password

## Démarrage rapide

### Créer votre premier client
1. Cliquer sur l'onglet "Créer un client"
2. Remplir nom et prénom (obligatoires)
3. Ajouter les autres informations
4. Cliquer "Enregistrer le client"

### Créer votre première facture
1. Aller dans "Toutes les factures"
2. Cliquer "Nouvelle facture"
3. Sélectionner un client
4. Remplir le montant HT
5. Cliquer "Enregistrer"

## Support rapide

### Problème de connexion à la base
```bash
sudo systemctl restart mysql
```

### Problème de port occupé
```bash
sudo lsof -i :8080
# Tuer le processus si nécessaire
```

### Réinitialiser la base de données
```bash
sudo mysql -e "DROP DATABASE factugest_db;"
sudo mysql < database/init.sql
```

