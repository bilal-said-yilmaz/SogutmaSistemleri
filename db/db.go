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
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			image VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS services (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			description TEXT,
			image VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS about (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			content TEXT,
			image VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS contact (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			phone VARCHAR(20),
			email VARCHAR(100),
			address TEXT
		)`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Admin kullanıcısını oluştur (eğer yoksa)
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", "admin").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		_, err = DB.Exec(
			"INSERT INTO users (username, password) VALUES ($1, $2)",
			"admin",
			"admin123", // Gerçek uygulamada şifre hash'lenmelidir!
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
