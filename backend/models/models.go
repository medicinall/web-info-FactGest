package models

import (
	"time"
)

// Client représente un client de l'entreprise
type Client struct {
	ID          int       `json:"id" db:"id"`
	Nom         string    `json:"nom" db:"nom"`
	Prenom      string    `json:"prenom" db:"prenom"`
	Adresse     *string   `json:"adresse" db:"adresse"`
	Ville       *string   `json:"ville" db:"ville"`
	CodePostal  *string   `json:"code_postal" db:"code_postal"`
	Telephone   *string   `json:"telephone" db:"telephone"`
	Email       *string   `json:"email" db:"email"`
	Entreprise  *string   `json:"entreprise" db:"entreprise"`
	SIRET       *string   `json:"siret" db:"siret"`
	DateCreation time.Time `json:"date_creation" db:"date_creation"`
	DateModification time.Time `json:"date_modification" db:"date_modification"`
}

// Facture représente une facture
type Facture struct {
	ID              int       `json:"id" db:"id"`
	NumeroFacture   string    `json:"numero_facture" db:"numero_facture"`
	ClientID        int       `json:"client_id" db:"client_id"`
	DateFacture     time.Time `json:"date_facture" db:"date_facture"`
	DateEcheance    time.Time `json:"date_echeance" db:"date_echeance"`
	Statut          string    `json:"statut" db:"statut"` // "brouillon", "envoyee", "payee", "en_retard"
	MontantHT       float64   `json:"montant_ht" db:"montant_ht"`
	TauxTVA         float64   `json:"taux_tva" db:"taux_tva"`
	MontantTVA      float64   `json:"montant_tva" db:"montant_tva"`
	MontantTTC      float64   `json:"montant_ttc" db:"montant_ttc"`
	Description     *string   `json:"description" db:"description"`
	Notes           *string   `json:"notes" db:"notes"`
	DateCreation    time.Time `json:"date_creation" db:"date_creation"`
	DateModification time.Time `json:"date_modification" db:"date_modification"`
	
	// Relation avec le client
	Client *Client `json:"client,omitempty"`
}

// LigneFacture représente une ligne d'une facture
type LigneFacture struct {
	ID          int     `json:"id" db:"id"`
	FactureID   int     `json:"facture_id" db:"facture_id"`
	Description string  `json:"description" db:"description"`
	Quantite    float64 `json:"quantite" db:"quantite"`
	PrixUnitaire float64 `json:"prix_unitaire" db:"prix_unitaire"`
	MontantHT   float64 `json:"montant_ht" db:"montant_ht"`
}

// Utilisateur pour l'authentification basique
type Utilisateur struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"` // Le mot de passe ne sera jamais retourné en JSON
	Role     string `json:"role" db:"role"`  // "admin", "user"
	Actif    bool   `json:"actif" db:"actif"`
}

