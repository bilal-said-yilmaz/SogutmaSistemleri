package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database")

	// Create tables
	createTables()
}

func createTables() {
	// Products table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2),
			image TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Users table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role_id INTEGER NOT NULL DEFAULT 2,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Services table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS services (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			image TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// About table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS about (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			image TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Contact table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS contact (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			phone VARCHAR(50),
			email VARCHAR(255),
			address TEXT,
			latitude DECIMAL(10,8),
			longitude DECIMAL(11,8),
			weekday_hours VARCHAR(100),
			saturday_hours VARCHAR(100),
			sunday_hours VARCHAR(100)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Hero table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS hero (
			id SERIAL PRIMARY KEY,
			subheading VARCHAR(255),
			heading VARCHAR(255),
			button_text VARCHAR(100),
			background_image TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Footer table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS footer (
			id SERIAL PRIMARY KEY,
			copyright TEXT,
			social_links JSONB,
			links JSONB
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully created tables")
}
