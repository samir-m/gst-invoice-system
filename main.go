package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var tmpl *template.Template

func init() {
	var err error

	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("ERROR parsing templates:", err)
	}

	tmpl, err = tmpl.ParseGlob("templates/pages/*.html")
	if err != nil {
		log.Fatal("ERROR parsing pages:", err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	db_host := os.Getenv("DB_HOST")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_database := os.Getenv("DB_DATABASE")

	dbConfig := DBConfig{
		host:     db_host,
		user:     db_user,
		password: db_password,
		database: db_database,
	}

	initDB(dbConfig)

	dir, _ := os.Getwd()
	fmt.Println("PWD:", dir)

	app := &App{
		Tmpl: tmpl,
	}

	r := NewRouter(app)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
