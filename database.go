package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {

	var err error

	db, err = sql.Open("sqlite", "invoice.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {

	queries := []string{
		`CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			phone TEXT,
			gstin TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product TEXT,
			price REAL,
			gst REAL
		);`,

		`CREATE TABLE IF NOT EXISTS invoices (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	invoice_no TEXT UNIQUE,
	customer_id INTEGER,
	subtotal REAL,
	gst_total REAL,
	grand_total REAL,
	created_at TEXT
);`,

		`CREATE TABLE IF NOT EXISTS invoice_items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	invoice_id INTEGER,
	product_id INTEGER,
	qty INTEGER,
	price REAL,
	gst REAL,
	line_total REAL
);`,
		`CREATE TABLE IF NOT EXISTS company (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_name TEXT NOT NULL,
    gstin TEXT,
    phone TEXT,
    email TEXT,
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
