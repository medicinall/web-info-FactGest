package main

import (
	"log"
	"net/http"

	"factugest-webinformatique/config"
	"factugest-webinformatique/database"
	"factugest-webinformatique/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Charger la configuration
	cfg := config.LoadConfig()

	// Initialiser la base de données
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Erreur lors de l'initialisation de la base de données:", err)
	}
	defer db.Close()

	// Créer le routeur
	router := mux.NewRouter()

	// Initialiser les handlers
	h := handlers.NewHandlers(db)

	// Routes API
	api := router.PathPrefix("/api").Subrouter()
	
	// Routes pour les clients
	api.HandleFunc("/clients", h.GetClients).Methods("GET")
	api.HandleFunc("/clients", h.CreateClient).Methods("POST")
	api.HandleFunc("/clients/{id}", h.GetClient).Methods("GET")
	api.HandleFunc("/clients/{id}", h.UpdateClient).Methods("PUT")
	api.HandleFunc("/clients/{id}", h.DeleteClient).Methods("DELETE")

	// Routes pour les factures
	api.HandleFunc("/factures", h.GetFactures).Methods("GET")
	api.HandleFunc("/factures", h.CreateFacture).Methods("POST")
	api.HandleFunc("/factures/{id}", h.GetFacture).Methods("GET")
	api.HandleFunc("/factures/{id}", h.UpdateFacture).Methods("PUT")
	api.HandleFunc("/factures/{id}", h.DeleteFacture).Methods("DELETE")
	api.HandleFunc("/factures/{id}/pdf", h.GenerateFacturePDF).Methods("GET")

	// Routes pour l'authentification
	api.HandleFunc("/login", h.Login).Methods("POST")
	api.HandleFunc("/logout", h.Logout).Methods("POST")

	// Routes pour l'administration
	api.HandleFunc("/admin/stats", h.GetStats).Methods("GET")
	api.HandleFunc("/admin/utilisateurs", h.GetUtilisateurs).Methods("GET")

	// Servir les fichiers statiques
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/")))

	// Configuration CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Printf("Serveur démarré sur le port %s", cfg.Port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.Port, handler))
}

