package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"factugest-webinformatique/models"

	"github.com/gorilla/mux"
)

// ===== HANDLERS POUR LES FACTURES =====

// GetFactures récupère toutes les factures avec les informations client
func (h *Handlers) GetFactures(w http.ResponseWriter, r *http.Request) {
	query := `SELECT f.id, f.numero_facture, f.client_id, f.date_facture, f.date_echeance, 
			  f.statut, f.montant_ht, f.taux_tva, f.montant_tva, f.montant_ttc, 
			  f.description, f.notes, f.date_creation, f.date_modification,
			  c.nom, c.prenom, c.entreprise
			  FROM factures f 
			  LEFT JOIN clients c ON f.client_id = c.id 
			  ORDER BY f.date_facture DESC`
	
	rows, err := h.DB.Query(query)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la récupération des factures",
		})
		return
	}
	defer rows.Close()

	var factures []models.Facture
	for rows.Next() {
		var facture models.Facture
		var client models.Client
		
		err := rows.Scan(&facture.ID, &facture.NumeroFacture, &facture.ClientID, 
			&facture.DateFacture, &facture.DateEcheance, &facture.Statut, 
			&facture.MontantHT, &facture.TauxTVA, &facture.MontantTVA, &facture.MontantTTC,
			&facture.Description, &facture.Notes, &facture.DateCreation, &facture.DateModification,
			&client.Nom, &client.Prenom, &client.Entreprise)
		
		if err != nil {
			h.sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Erreur lors du scan des factures",
			})
			return
		}
		
		client.ID = facture.ClientID
		facture.Client = &client
		factures = append(factures, facture)
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    factures,
	})
}

// GetFacture récupère une facture par son ID
func (h *Handlers) GetFacture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID facture invalide",
		})
		return
	}

	query := `SELECT f.id, f.numero_facture, f.client_id, f.date_facture, f.date_echeance, 
			  f.statut, f.montant_ht, f.taux_tva, f.montant_tva, f.montant_ttc, 
			  f.description, f.notes, f.date_creation, f.date_modification,
			  c.nom, c.prenom, c.adresse, c.ville, c.code_postal, c.telephone, c.email, c.entreprise, c.siret
			  FROM factures f 
			  LEFT JOIN clients c ON f.client_id = c.id 
			  WHERE f.id = ?`
	
	var facture models.Facture
	var client models.Client
	
	err = h.DB.QueryRow(query, id).Scan(&facture.ID, &facture.NumeroFacture, &facture.ClientID, 
		&facture.DateFacture, &facture.DateEcheance, &facture.Statut, 
		&facture.MontantHT, &facture.TauxTVA, &facture.MontantTVA, &facture.MontantTTC,
		&facture.Description, &facture.Notes, &facture.DateCreation, &facture.DateModification,
		&client.Nom, &client.Prenom, &client.Adresse, &client.Ville, &client.CodePostal,
		&client.Telephone, &client.Email, &client.Entreprise, &client.SIRET)
	
	if err != nil {
		h.sendJSONResponse(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Facture non trouvée",
		})
		return
	}

	client.ID = facture.ClientID
	facture.Client = &client

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    facture,
	})
}

// CreateFacture crée une nouvelle facture
func (h *Handlers) CreateFacture(w http.ResponseWriter, r *http.Request) {
	var facture models.Facture
	if err := json.NewDecoder(r.Body).Decode(&facture); err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Données JSON invalides",
		})
		return
	}

	// Validation des champs obligatoires
	if facture.ClientID == 0 || facture.NumeroFacture == "" {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Le client et le numéro de facture sont obligatoires",
		})
		return
	}

	// Vérifier que le client existe
	var clientExists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clients WHERE id = ?)", facture.ClientID).Scan(&clientExists)
	if err != nil || !clientExists {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Client inexistant",
		})
		return
	}

	// Calculer les montants
	facture.MontantTVA = facture.MontantHT * facture.TauxTVA / 100
	facture.MontantTTC = facture.MontantHT + facture.MontantTVA

	query := `INSERT INTO factures (numero_facture, client_id, date_facture, date_echeance, 
			  statut, montant_ht, taux_tva, montant_tva, montant_ttc, description, notes) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := h.DB.Exec(query, facture.NumeroFacture, facture.ClientID, 
		facture.DateFacture, facture.DateEcheance, facture.Statut, 
		facture.MontantHT, facture.TauxTVA, facture.MontantTVA, facture.MontantTTC,
		facture.Description, facture.Notes)
	
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la création de la facture",
		})
		return
	}

	id, _ := result.LastInsertId()
	facture.ID = int(id)

	h.sendJSONResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "Facture créée avec succès",
		Data:    facture,
	})
}

// UpdateFacture met à jour une facture existante
func (h *Handlers) UpdateFacture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID facture invalide",
		})
		return
	}

	var facture models.Facture
	if err := json.NewDecoder(r.Body).Decode(&facture); err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Données JSON invalides",
		})
		return
	}

	// Validation des champs obligatoires
	if facture.ClientID == 0 || facture.NumeroFacture == "" {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Le client et le numéro de facture sont obligatoires",
		})
		return
	}

	// Calculer les montants
	facture.MontantTVA = facture.MontantHT * facture.TauxTVA / 100
	facture.MontantTTC = facture.MontantHT + facture.MontantTVA

	query := `UPDATE factures SET numero_facture = ?, client_id = ?, date_facture = ?, 
			  date_echeance = ?, statut = ?, montant_ht = ?, taux_tva = ?, montant_tva = ?, 
			  montant_ttc = ?, description = ?, notes = ?, date_modification = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	
	result, err := h.DB.Exec(query, facture.NumeroFacture, facture.ClientID, 
		facture.DateFacture, facture.DateEcheance, facture.Statut, 
		facture.MontantHT, facture.TauxTVA, facture.MontantTVA, facture.MontantTTC,
		facture.Description, facture.Notes, id)
	
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la mise à jour de la facture",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.sendJSONResponse(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Facture non trouvée",
		})
		return
	}

	facture.ID = id
	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Facture mise à jour avec succès",
		Data:    facture,
	})
}

// DeleteFacture supprime une facture
func (h *Handlers) DeleteFacture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID facture invalide",
		})
		return
	}

	query := `DELETE FROM factures WHERE id = ?`
	result, err := h.DB.Exec(query, id)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la suppression de la facture",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.sendJSONResponse(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Facture non trouvée",
		})
		return
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Facture supprimée avec succès",
	})
}



// GenerateNumeroFacture génère un numéro de facture automatique
func (h *Handlers) GenerateNumeroFacture() string {
	year := time.Now().Year()
	
	// Compter le nombre de factures de l'année en cours
	var count int
	query := `SELECT COUNT(*) FROM factures WHERE YEAR(date_facture) = ?`
	h.DB.QueryRow(query, year).Scan(&count)
	
	return fmt.Sprintf("FACT-%d-%04d", year, count+1)
}

