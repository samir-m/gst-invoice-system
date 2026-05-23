package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB(config DBConfig) {

	dsn := config.user + ":" + config.password + "@tcp(" + config.host + ")/" + config.database + "?parseTime=true"

	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// IMPORTANT: verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("database unreachable: %v", err)
	}

	createTables()
}

func createTables() {

	queries := []string{
		`CREATE TABLE IF NOT EXISTS customers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(50),
    gstin VARCHAR(100)
);`,

		`CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product VARCHAR(255),
    price DECIMAL(10,2),
    gst DECIMAL(5,2)
);`,

		`CREATE TABLE IF NOT EXISTS invoices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    invoice_no VARCHAR(100) UNIQUE,
    customer_id INT,
    subtotal DECIMAL(10,2),
    gst_total DECIMAL(10,2),
    grand_total DECIMAL(10,2),
    created_at DATETIME
);`,

		`CREATE TABLE IF NOT EXISTS invoice_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    invoice_id INT,
    product_id INT,
    qty INT,
    price DECIMAL(10,2),
    gst DECIMAL(5,2),
    line_total DECIMAL(10,2)
);`,
		`CREATE TABLE IF NOT EXISTS company (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    gstin VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal(err)
		}
	}
}
