package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"factugest-webinformatique/pdf"

	"github.com/gorilla/mux"
)

// GenerateFacturePDF génère et retourne un PDF de facture
func (h *Handlers) GenerateFacturePDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID facture invalide",
		})
		return
	}

	// Récupérer la facture
	facture, err := h.getFactureByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.sendJSONResponse(w, http.StatusNotFound, Response{
				Success: false,
				Error:   "Facture non trouvée",
			})
		} else {
			h.sendJSONResponse(w, http.StatusInternalServerError, Response{
				Success: false,
				Error:   "Erreur lors de la récupération de la facture: " + err.Error(),
			})
		}
		return
	}

	// Récupérer le client associé
	client, err := h.getClientByIDForPDF(facture.ClientID)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la récupération du client: " + err.Error(),
		})
		return
	}

	// Générer le PDF
	generator := pdf.NewPDFGenerator("assets/logo.png") // Logo optionnel
	pdfBuffer, err := generator.GenerateFacturePDF(facture, client)
	if err != nil {
		h.sendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Erreur lors de la génération du PDF: " + err.Error(),
		})
		return
	}

	// Définir les en-têtes pour le téléchargement du PDF
	filename := fmt.Sprintf("facture_%s.pdf", facture.NumeroFacture)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(pdfBuffer.Len()))

	// Envoyer le PDF
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuffer.Bytes())
}

// getFactureByID récupère une facture par son ID
func (h *Handlers) getFactureByID(id int) (*pdf.FacturePDF, error) {
	query := `SELECT id, numero_facture, client_id, date_facture, date_echeance, statut, 
			  montant_ht, taux_tva, montant_tva, montant_ttc, 
			  COALESCE(description, '') as description,
			  COALESCE(notes, '') as notes,
			  date_creation, date_modification 
			  FROM factures WHERE id = ?`
	
	var facture pdf.FacturePDF
	var description, notes string
	
	err := h.DB.QueryRow(query, id).Scan(
		&facture.ID, &facture.NumeroFacture, &facture.ClientID,
		&facture.DateFacture, &facture.DateEcheance, &facture.Statut,
		&facture.MontantHT, &facture.TauxTVA, &facture.MontantTVA, &facture.MontantTTC,
		&description, &notes,
		&facture.DateCreation, &facture.DateModification,
	)
	
	if err != nil {
		return nil, err
	}

	// Convertir les strings en pointeurs si non vides
	if description != "" {
		facture.Description = &description
	}
	if notes != "" {
		facture.Notes = &notes
	}

	return &facture, nil
}

// getClientByIDForPDF récupère un client par son ID pour le PDF
func (h *Handlers) getClientByIDForPDF(id int) (*pdf.ClientSimple, error) {
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
	
	var client pdf.ClientSimple
	var adresse, ville, codePostal, telephone, email, entreprise, siret string
	
	err := h.DB.QueryRow(query, id).Scan(
		&client.ID, &client.Nom, &client.Prenom,
		&adresse, &ville, &codePostal, &telephone, &email, &entreprise, &siret,
		&client.DateCreation, &client.DateModification,
	)
	
	if err != nil {
		return nil, err
	}

	// Convertir les strings en pointeurs si non vides
	if adresse != "" {
		client.Adresse = &adresse
	}
	if ville != "" {
		client.Ville = &ville
	}
	if codePostal != "" {
		client.CodePostal = &codePostal
	}
	if telephone != "" {
		client.Telephone = &telephone
	}
	if email != "" {
		client.Email = &email
	}
	if entreprise != "" {
		client.Entreprise = &entreprise
	}
	if siret != "" {
		client.SIRET = &siret
	}

	return &client, nil
}

