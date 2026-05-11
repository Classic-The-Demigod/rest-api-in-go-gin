package main

import (
	"database/sql"
	"log"

	_ "rest-api-in-go/docs"
	"rest-api-in-go/internal/database"
	"rest-api-in-go/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// @title Go Gin Rest API
// @version 1.0
// @description A rest API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**
// @name Authorization

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
