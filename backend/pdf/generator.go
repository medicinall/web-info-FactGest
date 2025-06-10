package pdf

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct {
	LogoPath string
}

func NewPDFGenerator(logoPath string) *PDFGenerator {
	return &PDFGenerator{LogoPath: logoPath}
}

// ClientSimple pour la génération PDF
type ClientSimple struct {
	ID          int       `json:"id"`
	Nom         string    `json:"nom"`
	Prenom      string    `json:"prenom"`
	Adresse     *string   `json:"adresse"`
	Ville       *string   `json:"ville"`
	CodePostal  *string   `json:"code_postal"`
	Telephone   *string   `json:"telephone"`
	Email       *string   `json:"email"`
	Entreprise  *string   `json:"entreprise"`
	SIRET       *string   `json:"siret"`
	DateCreation time.Time `json:"date_creation"`
	DateModification time.Time `json:"date_modification"`
}

// Facture simplifiée pour la génération PDF
type FacturePDF struct {
	ID              int       `json:"id"`
	NumeroFacture   string    `json:"numero_facture"`
	ClientID        int       `json:"client_id"`
	DateFacture     time.Time `json:"date_facture"`
	DateEcheance    time.Time `json:"date_echeance"`
	Statut          string    `json:"statut"`
	MontantHT       float64   `json:"montant_ht"`
	TauxTVA         float64   `json:"taux_tva"`
	MontantTVA      float64   `json:"montant_tva"`
	MontantTTC      float64   `json:"montant_ttc"`
	Description     *string   `json:"description"`
	Notes           *string   `json:"notes"`
	DateCreation    time.Time `json:"date_creation"`
	DateModification time.Time `json:"date_modification"`
}

// GenerateFacturePDF génère un PDF de facture basé sur le modèle WebInformatique
func (g *PDFGenerator) GenerateFacturePDF(facture *FacturePDF, client *ClientSimple) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)

	pdf.AddPage()

	// Logo (si disponible)
	if g.LogoPath != "" {
		if _, err := os.Stat(g.LogoPath); err == nil {
			pdf.Image(g.LogoPath, 20, 20, 60, 20, false, "", 0, "")
		}
	}

	// Position après le logo
	y := 45.0

	// Titre de la facture
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "FACTURE", "0", 1, "C", false, 0, "")
	y += 15

	// Informations de la facture (en-tête)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("Facture N° %s", facture.NumeroFacture), "0", 1, "C", false, 0, "")
	y += 10

	// Tableau d'informations principales
	pdf.SetFont("Arial", "", 10)
	infoHeader := []string{"Date de facture", "Date d'échéance", "Statut"}
	infoData := [][]string{
		infoHeader,
		{
			facture.DateFacture.Format("02/01/2006"),
			facture.DateEcheance.Format("02/01/2006"),
			getStatutLabel(facture.Statut),
		},
	}

	drawTable(pdf, 20, y, []float64{63, 63, 64}, infoData)
	y += 25

	// Section Expéditeur et Destinataire
	expediteurData := [][]string{
		{"EXPÉDITEUR", "", "DESTINATAIRE"},
		{"WEB INFORMATIQUE SARL", "", client.Nom + " " + client.Prenom},
		{"154 bis rue du général de Gaulle", "", getStringValue(client.Entreprise)},
		{"76770 LE HOULME", "", getStringValue(client.Adresse)},
		{"Tél. 06.99.50.76.76", "", getStringValue(client.CodePostal) + " " + getStringValue(client.Ville)},
		{"Tél. 02.35.74.19.29", "", getStringValue(client.Telephone)},
		{"contact@webinformatique.eu", "", getStringValue(client.Email)},
		{"SIRET: 493 928 139 00010", "", "SIRET: " + getStringValue(client.SIRET)},
	}

	drawTable(pdf, 20, y, []float64{95, 0, 95}, expediteurData)
	y += 55

	// Description des prestations
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, "DESCRIPTION DES PRESTATIONS", "0", 1, "L", false, 0, "")
	y += 12

	// Tableau des prestations
	prestationHeader := []string{"Désignation", "Quantité", "Prix unitaire HT", "Montant HT"}
	prestationData := [][]string{prestationHeader}

	// Ajouter la description de la facture comme ligne de prestation
	description := getStringValue(facture.Description)
	if description == "" {
		description = "Prestation de service informatique"
	}

	prestationData = append(prestationData, []string{
		description,
		"1",
		fmt.Sprintf("%.2f €", facture.MontantHT),
		fmt.Sprintf("%.2f €", facture.MontantHT),
	})

	// Ajouter des lignes vides pour l'esthétique
	for i := 0; i < 3; i++ {
		prestationData = append(prestationData, []string{"", "", "", ""})
	}

	drawTable(pdf, 20, y, []float64{95, 30, 35, 30}, prestationData)
	y += 35

	// Totaux
	totauxData := [][]string{
		{"", "", "TOTAL HT", fmt.Sprintf("%.2f €", facture.MontantHT)},
		{"", "", fmt.Sprintf("TVA (%.1f%%)", facture.TauxTVA), fmt.Sprintf("%.2f €", facture.MontantTVA)},
		{"", "", "TOTAL TTC", fmt.Sprintf("%.2f €", facture.MontantTTC)},
	}

	pdf.SetFont("Arial", "B", 10)
	drawTable(pdf, 20, y, []float64{95, 30, 35, 30}, totauxData)
	y += 25

	// Notes (si présentes)
	notes := getStringValue(facture.Notes)
	if notes != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(190, 8, "NOTES :", "0", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 9)
		pdf.MultiCell(190, 5, notes, "1", "L", false)
		y += 20
	}

	// Conditions de paiement
	y += 10
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(190, 8, "CONDITIONS DE PAIEMENT", "0", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	conditions := `Paiement à réception de facture.
ATTENTION : LES RÈGLEMENTS PAR CHÈQUES NE SONT PAS ACCEPTÉS.
Tout retard de paiement entraînera l'application d'intérêts de retard au taux légal en vigueur.`
	
	pdf.MultiCell(190, 5, conditions, "1", "L", false)
	y += 25

	// Pied de page avec mentions légales
	pdf.SetY(260) // Position fixe en bas de page
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(190, 5, "Web Informatique SARL au capital de 8.000€", "0", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, "SIRET: 493 928 139 00010 - N° TVA intracommunautaire : FR493928139", "0", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, "154 bis rue du général de Gaulle - 76770 LE HOULME", "0", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// drawTable dessine un tableau avec les données fournies
func drawTable(pdf *gofpdf.Fpdf, x, y float64, widths []float64, data [][]string) {
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.2)

	for i, row := range data {
		for j, cell := range row {
			if j < len(widths) && widths[j] > 0 { // Ignorer les colonnes avec largeur 0
				pdf.SetXY(x+sum(widths[:j]), y+float64(i*8))
				
				// Style différent pour l'en-tête
				if i == 0 {
					pdf.SetFont("Arial", "B", 9)
					pdf.CellFormat(widths[j], 8, cell, "1", 0, "C", true, 0, "")
				} else {
					pdf.SetFont("Arial", "", 9)
					pdf.CellFormat(widths[j], 8, cell, "1", 0, "C", false, 0, "")
				}
			}
		}
	}
}

// sum calcule la somme d'un slice de float64
func sum(nums []float64) float64 {
	total := 0.0
	for _, num := range nums {
		total += num
	}
	return total
}

// getStatutLabel retourne le libellé français du statut
func getStatutLabel(statut string) string {
	switch statut {
	case "brouillon":
		return "Brouillon"
	case "envoyee":
		return "Envoyée"
	case "payee":
		return "Payée"
	case "en_retard":
		return "En retard"
	default:
		return "Inconnu"
	}
}

// getStringValue retourne la valeur d'un pointeur de string ou une chaîne vide
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

