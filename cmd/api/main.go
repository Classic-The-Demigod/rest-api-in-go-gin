package main

import (
	"database/sql"
	"log"
	"rest-api-in-go/internal/database"
	"rest-api-in-go/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 4000),
		jwtSecret: env.GetEnvString("JWT_SECRET", "supersecretkey"),
		models:    models,
	}
	if err := app.server(); err != nil {
		log.Fatal(err)
	}
}
