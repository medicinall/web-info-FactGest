-- Script d'initialisation de la base de données FactuGest-WebInformatique

-- Change le mot de passe root (pour localhost uniquement)
ALTER USER 'root'@'localhost' IDENTIFIED BY '';
FLUSH PRIVILEGES;

-- Créer la base de données
CREATE DATABASE IF NOT EXISTS factugest_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Créer l'utilisateur factuguest et lui donner les permissions
CREATE USER IF NOT EXISTS 'factuguest'@'localhost' IDENTIFIED BY 'FactuGest_P@ssw0rd!';
GRANT ALL PRIVILEGES ON factugest_db.* TO 'factuguest'@'localhost';
FLUSH PRIVILEGES;

-- Utiliser la base de données
USE factugest_db;

-- Table des utilisateurs
CREATE TABLE IF NOT EXISTS utilisateurs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'user') DEFAULT 'user',
    actif BOOLEAN DEFAULT TRUE,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table des clients
CREATE TABLE IF NOT EXISTS clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    prenom VARCHAR(100) NOT NULL,
    adresse TEXT,
    ville VARCHAR(100),
    code_postal VARCHAR(10),
    telephone VARCHAR(20),
    email VARCHAR(100),
    entreprise VARCHAR(150),
    siret VARCHAR(20),
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_modification TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_nom (nom),
    INDEX idx_email (email),
    INDEX idx_entreprise (entreprise)
);

-- Table des factures
CREATE TABLE IF NOT EXISTS factures (
    id INT AUTO_INCREMENT PRIMARY KEY,
    numero_facture VARCHAR(50) UNIQUE NOT NULL,
    client_id INT NOT NULL,
    date_facture DATE NOT NULL,
    date_echeance DATE NOT NULL,
    statut ENUM('brouillon', 'envoyee', 'payee', 'en_retard') DEFAULT 'brouillon',
    montant_ht DECIMAL(10,2) DEFAULT 0.00,
    taux_tva DECIMAL(5,2) DEFAULT 20.00,
    montant_tva DECIMAL(10,2) DEFAULT 0.00,
    montant_ttc DECIMAL(10,2) DEFAULT 0.00,
    description TEXT,
    notes TEXT,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date_modification TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    INDEX idx_numero_facture (numero_facture),
    INDEX idx_client_id (client_id),
    INDEX idx_date_facture (date_facture),
    INDEX idx_statut (statut)
);

-- Table des lignes de factures
CREATE TABLE IF NOT EXISTS lignes_factures (
    id INT AUTO_INCREMENT PRIMARY KEY,
    facture_id INT NOT NULL,
    description TEXT NOT NULL,
    quantite DECIMAL(10,2) NOT NULL DEFAULT 1.00,
    prix_unitaire DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    montant_ht DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    FOREIGN KEY (facture_id) REFERENCES factures(id) ON DELETE CASCADE,
    INDEX idx_facture_id (facture_id)
);

-- Insérer un utilisateur admin par défaut
INSERT INTO utilisateurs (username, password, role) VALUES 
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')
ON DUPLICATE KEY UPDATE username = username;

-- Insérer quelques clients de test
INSERT INTO clients (nom, prenom, adresse, ville, code_postal, telephone, email, entreprise) VALUES 
('Dupont', 'Jean', '123 Rue de la Paix', 'Paris', '75001', '01.23.45.67.89', 'jean.dupont@email.com', 'Dupont SARL'),
('Martin', 'Marie', '456 Avenue des Champs', 'Lyon', '69001', '04.56.78.90.12', 'marie.martin@email.com', 'Martin & Associés'),
('Bernard', 'Pierre', '789 Boulevard Saint-Germain', 'Marseille', '13001', '04.91.23.45.67', 'pierre.bernard@email.com', 'Bernard Consulting')
ON DUPLICATE KEY UPDATE nom = nom;
