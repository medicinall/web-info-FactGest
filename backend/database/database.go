package database

import (
	"database/sql"
	"fmt"

	"factugest-webinformatique/config"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture de la connexion à la base de données: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erreur lors de la connexion à la base de données: %v", err)
	}

	// Configuration de la pool de connexions
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

