package main

import "net/http"

func getProducts() ([]Product, error) {

	rows, err := db.Query("SELECT id, product, price, gst FROM products")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var c Product

		err := rows.Scan(&c.ID, &c.Product, &c.Price, &c.GST)
		if err != nil {
			return nil, err
		}

		products = append(products, c)
	}

	return products, nil
}

func (app *App) ProductHandler(w http.ResponseWriter, r *http.Request) {

	products, _ := getProducts()

	data := map[string]interface{}{
		"Title":    "Product",
		"Page":     "product",
		"Products": products,
	}

	app.Tmpl.ExecuteTemplate(w, "product", data)
}

func (app *App) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	product := r.FormValue("product")
	price := r.FormValue("price")
	gst := r.FormValue("gst")

	_, err := db.Exec(
		"INSERT INTO products(product, price, gst) VALUES(?,?, ?)",
		product,
		price,
		gst,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/product", 302)
}
