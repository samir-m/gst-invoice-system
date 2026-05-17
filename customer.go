package main

import (
	"net/http"
)

func getCustomers() ([]Customer, error) {

	rows, err := db.Query("SELECT id, name, phone, gstin FROM customers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var c Customer

		err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.GSTIN)
		if err != nil {
			return nil, err
		}

		customers = append(customers, c)
	}

	return customers, nil
}

func (app *App) CustomerHandler(w http.ResponseWriter, r *http.Request) {

	customers, _ := getCustomers()

	data := map[string]interface{}{
		"Title":     "Customer",
		"Page":      "customer",
		"Customers": customers,
	}

	app.Tmpl.ExecuteTemplate(w, "customer", data)
}

func (app *App) AddCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	name := r.FormValue("customer_name")
	phone := r.FormValue("phone_number")
	gstin := r.FormValue("gstin")

	_, err := db.Exec(
		"INSERT INTO customers(name, phone, gstin) VALUES(?,?, ?)",
		name,
		phone,
		gstin,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/customer", 302)
}
