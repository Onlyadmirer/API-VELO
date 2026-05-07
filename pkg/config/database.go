package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConnectDB membuka koneksi ke SQL database (PostgreSQL/MySQL) menggunakan string koneksi dsn.
// Fungsi ini mengembalikan *sql.DB yang bisa di-inject ke berbagai Repository.
func ConnectDB() (*sql.DB, error) {

	dbUrl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("database connection successfully")

	return db, nil

}
