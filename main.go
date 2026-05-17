package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

	initDB()

	dir, _ := os.Getwd()
	fmt.Println("PWD:", dir)

	app := &App{
		Tmpl: tmpl,
	}

	r := NewRouter(app)

	log.Fatal(http.ListenAndServe(":8080", r))
}
