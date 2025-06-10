package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handlers struct {
	DB *sql.DB
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

// Response structure pour les réponses JSON
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Client simplifié pour éviter les problèmes de NULL
type ClientSimple struct {
	ID          int       `json:"id"`
	Nom         string    `json:"nom"`
	Prenom      string    `json:"prenom"`
	Adresse     string    `json:"adresse"`
	Ville       string    `json:"ville"`
	CodePostal  string    `json:"code_postal"`
	Telephone   string    `json:"telephone"`
	Email       string    `json:"email"`
	Entreprise  string    `json:"entreprise"`
	SIRET       string    `json:"siret"`
	DateCreation time.Time `json:"date_creation"`
	DateModification time.Time `json:"date_modification"`
}

// Fonction utilitaire pour envoyer une réponse JSON
func (h *Handlers) sendJSONResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// GetClients récupère tous les clients
func (h *Handlers) GetClients(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, nom, prenom, 
			  COALESCE(adresse, '') as adresse,
			  COALESCE(ville, '') as ville,
			  COALESCE(code_postal, '') as code_postal,
			  COALESCE(telephone, '') as telephone,
			  COALESCE(email, '') as email,
			  COALESCE(entreprise, '') as entreprise,
			  COALESCE(siret, '') as siret,
			  date_creation, date_modification 
			  FROM clients ORDER BY nom, prenom`
	
	rows, err := h.DB.Query(query)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la récupération des clients: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var clients []ClientSimple
	for rows.Next() {
		var client ClientSimple
		err := rows.Scan(&client.ID, &client.Nom, &client.Prenom, &client.Adresse, 
			&client.Ville, &client.CodePostal, &client.Telephone, &client.Email, 
			&client.Entreprise, &client.SIRET, &client.DateCreation, &client.DateModification)
		if err != nil {
			h.sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Erreur lors du scan des clients: " + err.Error(),
			})
			return
		}
		clients = append(clients, client)
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    clients,
	})
}

// GetClient récupère un client par son ID
func (h *Handlers) GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID client invalide",
		})
		return
	}

	query := `SELECT id, nom, prenom, 
			  COALESCE(adresse, '') as adresse,
			  COALESCE(ville, '') as ville,
			  COALESCE(code_postal, '') as code_postal,
			  COALESCE(telephone, '') as telephone,
			  COALESCE(email, '') as email,
			  COALESCE(entreprise, '') as entreprise,
			  COALESCE(siret, '') as siret,
			  date_creation, date_modification 
			  FROM clients WHERE id = ?`
	
	var client ClientSimple
	err = h.DB.QueryRow(query, id).Scan(&client.ID, &client.Nom, &client.Prenom, 
		&client.Adresse, &client.Ville, &client.CodePostal, &client.Telephone, 
		&client.Email, &client.Entreprise, &client.SIRET, &client.DateCreation, &client.DateModification)
	
	if err == sql.ErrNoRows {
		h.sendJSONResponse(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Client non trouvé",
		})
		return
	} else if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la récupération du client: " + err.Error(),
		})
		return
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    client,
	})
}

// CreateClient crée un nouveau client
func (h *Handlers) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client ClientSimple
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Données JSON invalides",
		})
		return
	}

	// Validation des champs obligatoires
	if client.Nom == "" || client.Prenom == "" {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Le nom et le prénom sont obligatoires",
		})
		return
	}

	query := `INSERT INTO clients (nom, prenom, adresse, ville, code_postal, telephone, email, entreprise, siret) 
			  VALUES (?, ?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''))`
	
	result, err := h.DB.Exec(query, client.Nom, client.Prenom, client.Adresse, 
		client.Ville, client.CodePostal, client.Telephone, client.Email, 
		client.Entreprise, client.SIRET)
	
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la création du client: " + err.Error(),
		})
		return
	}

	id, _ := result.LastInsertId()
	client.ID = int(id)

	h.sendJSONResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "Client créé avec succès",
		Data:    client,
	})
}

// UpdateClient met à jour un client existant
func (h *Handlers) UpdateClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID client invalide",
		})
		return
	}

	var client ClientSimple
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Données JSON invalides",
		})
		return
	}

	// Validation des champs obligatoires
	if client.Nom == "" || client.Prenom == "" {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Le nom et le prénom sont obligatoires",
		})
		return
	}

	query := `UPDATE clients SET nom = ?, prenom = ?, adresse = NULLIF(?, ''), ville = NULLIF(?, ''), 
			  code_postal = NULLIF(?, ''), telephone = NULLIF(?, ''), email = NULLIF(?, ''), 
			  entreprise = NULLIF(?, ''), siret = NULLIF(?, ''), date_modification = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	
	result, err := h.DB.Exec(query, client.Nom, client.Prenom, client.Adresse, 
		client.Ville, client.CodePostal, client.Telephone, client.Email, 
		client.Entreprise, client.SIRET, id)
	
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la mise à jour du client: " + err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.sendJSONResponse(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Client non trouvé",
		})
		return
	}

	client.ID = id
	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Client mis à jour avec succès",
		Data:    client,
	})
}

// DeleteClient supprime un client
func (h *Handlers) DeleteClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID client invalide",
		})
		return
	}

	// Vérifier s'il y a des factures associées à ce client
	var count int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM factures WHERE client_id = ?", id).Scan(&count)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la vérification des factures",
		})
		return
	}

	if count > 0 {
		h.sendJSONResponse(w, http.StatusConflict, Response{
			Success: false,
			Error:   fmt.Sprintf("Impossible de supprimer le client : %d facture(s) associée(s)", count),
		})
		return
	}

	query := `DELETE FROM clients WHERE id = ?`
	result, err := h.DB.Exec(query, id)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la suppression du client: " + err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.sendJSONResponse(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Client non trouvé",
		})
		return
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Client supprimé avec succès",
	})
}

