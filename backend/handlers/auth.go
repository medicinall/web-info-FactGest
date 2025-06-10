package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ===== HANDLERS POUR L'AUTHENTIFICATION =====

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	} `json:"user"`
}

// Login authentifie un utilisateur
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Données JSON invalides",
		})
		return
	}

	// Récupérer l'utilisateur de la base de données
	var userID int
	var username, hashedPassword, role string
	var actif bool
	
	query := `SELECT id, username, password, role, actif FROM utilisateurs WHERE username = ?`
	err := h.DB.QueryRow(query, loginReq.Username).Scan(&userID, &username, &hashedPassword, &role, &actif)
	
	if err != nil {
		h.sendJSONResponse(w, http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Nom d'utilisateur ou mot de passe incorrect",
		})
		return
	}

	if !actif {
		h.sendJSONResponse(w, http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Compte désactivé",
		})
		return
	}

	// Vérifier le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password))
	if err != nil {
		h.sendJSONResponse(w, http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Nom d'utilisateur ou mot de passe incorrect",
		})
		return
	}

	// Créer le token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token valide 24h
	})

	tokenString, err := token.SignedString([]byte("factugest-secret-key-2024"))
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la génération du token",
		})
		return
	}

	// Préparer la réponse
	loginResp := LoginResponse{
		Token: tokenString,
	}
	loginResp.User.ID = userID
	loginResp.User.Username = username
	loginResp.User.Role = role

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Connexion réussie",
		Data:    loginResp,
	})
}

// Logout déconnecte un utilisateur (côté client principalement)
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Déconnexion réussie",
	})
}

// ===== HANDLERS POUR L'ADMINISTRATION =====

// GetStats retourne les statistiques générales
func (h *Handlers) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := make(map[string]interface{})

	// Nombre total de clients
	var totalClients int
	h.DB.QueryRow("SELECT COUNT(*) FROM clients").Scan(&totalClients)
	stats["total_clients"] = totalClients

	// Nombre total de factures
	var totalFactures int
	h.DB.QueryRow("SELECT COUNT(*) FROM factures").Scan(&totalFactures)
	stats["total_factures"] = totalFactures

	// Montant total des factures
	var montantTotal float64
	h.DB.QueryRow("SELECT COALESCE(SUM(montant_ttc), 0) FROM factures").Scan(&montantTotal)
	stats["montant_total"] = montantTotal

	// Factures par statut
	rows, err := h.DB.Query("SELECT statut, COUNT(*) FROM factures GROUP BY statut")
	if err == nil {
		facturesParStatut := make(map[string]int)
		for rows.Next() {
			var statut string
			var count int
			rows.Scan(&statut, &count)
			facturesParStatut[statut] = count
		}
		rows.Close()
		stats["factures_par_statut"] = facturesParStatut
	}

	// Factures récentes (5 dernières)
	recentQuery := `SELECT f.id, f.numero_facture, f.montant_ttc, f.statut, f.date_facture,
					c.nom, c.prenom, c.entreprise
					FROM factures f 
					LEFT JOIN clients c ON f.client_id = c.id 
					ORDER BY f.date_creation DESC LIMIT 5`
	
	rows, err = h.DB.Query(recentQuery)
	if err == nil {
		var facturesRecentes []map[string]interface{}
		for rows.Next() {
			var id int
			var numeroFacture string
			var montantTTC float64
			var statut string
			var dateFacture time.Time
			var nom, prenom, entreprise string
			
			rows.Scan(&id, &numeroFacture, &montantTTC, &statut, &dateFacture, &nom, &prenom, &entreprise)
			
			facture := map[string]interface{}{
				"id":             id,
				"numero_facture": numeroFacture,
				"montant_ttc":    montantTTC,
				"statut":         statut,
				"date_facture":   dateFacture,
				"client_nom":     nom + " " + prenom,
				"client_entreprise": entreprise,
			}
			facturesRecentes = append(facturesRecentes, facture)
		}
		rows.Close()
		stats["factures_recentes"] = facturesRecentes
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    stats,
	})
}

// GetUtilisateurs retourne la liste des utilisateurs (admin seulement)
func (h *Handlers) GetUtilisateurs(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, username, role, actif, date_creation FROM utilisateurs ORDER BY username`
	
	rows, err := h.DB.Query(query)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la récupération des utilisateurs",
		})
		return
	}
	defer rows.Close()

	var utilisateurs []map[string]interface{}
	for rows.Next() {
		var id int
		var username, role string
		var actif bool
		var dateCreation time.Time
		
		err := rows.Scan(&id, &username, &role, &actif, &dateCreation)
		if err != nil {
			continue
		}
		
		utilisateur := map[string]interface{}{
			"id":            id,
			"username":      username,
			"role":          role,
			"actif":         actif,
			"date_creation": dateCreation,
		}
		utilisateurs = append(utilisateurs, utilisateur)
	}

	h.sendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Data:    utilisateurs,
	})
}

