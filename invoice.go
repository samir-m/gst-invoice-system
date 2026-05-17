package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func getCompany() (Company, error) {
	var c Company

	query := `SELECT id, company_name, gstin, phone, email FROM company LIMIT 1`

	err := db.QueryRow(query).Scan(
		&c.ID,
		&c.CompanyName,
		&c.GSTIN,
		&c.Phone,
		&c.Email,
	)

	return c, err
}

func generateInvoiceNumber(companyName string) string {
	// Remove spaces and convert to uppercase
	name := strings.ToUpper(strings.ReplaceAll(companyName, " ", ""))

	// Take first 3 characters
	prefix := name
	if len(name) > 3 {
		prefix = name[:3]
	}

	// Date format: YYMMDD
	date := time.Now().Format("060102")

	// Random 4-digit number
	rand.Seed(time.Now().UnixNano())
	unique := rand.Intn(9000) + 1000

	// Final invoice number
	return fmt.Sprintf("%s-%s-%d", prefix, date, unique)
}

func getInvoices() ([]InvoiceList, error) {

	query := `
		SELECT 
			i.id,
			i.invoice_no,
			c.name,
			i.subtotal,
			i.gst_total,
			i.grand_total,
			i.created_at
		FROM invoices i
		JOIN customers c ON c.id = i.customer_id
		ORDER BY i.id DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	defer rows.Close()

	var invoices []InvoiceList

	for rows.Next() {
		var inv InvoiceList

		err := rows.Scan(
			&inv.ID,
			&inv.InvoiceNo,
			&inv.Customer,
			&inv.Subtotal,
			&inv.GSTTotal,
			&inv.GrandTotal,
			&inv.Date,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, inv)
	}

	return invoices, nil
}

func getInvoiceByID(id int) (InvoiceList, []InvoiceItemView, error) {

	// 1. Get invoice header
	query := `
		SELECT 
			i.id,
			i.invoice_no,
			c.name,
			i.subtotal,
			i.gst_total,
			i.grand_total,
			i.created_at
		FROM invoices i
		JOIN customers c ON c.id = i.customer_id
		WHERE i.id = ?
	`

	var inv InvoiceList

	err := db.QueryRow(query, id).Scan(
		&inv.ID,
		&inv.InvoiceNo,
		&inv.Customer,
		&inv.Subtotal,
		&inv.GSTTotal,
		&inv.GrandTotal,
		&inv.Date,
	)

	if err != nil {
		fmt.Println("err", err)
		return InvoiceList{}, nil, err
	}

	// 2. Get invoice items
	itemQuery := `
		SELECT 
			ii.id,
			ii.invoice_id,
			ii.product_id,
			p.product,
			ii.qty,
			ii.price,
			ii.gst,
			ii.line_total
		FROM invoice_items ii
		JOIN products p ON p.id = ii.product_id
		WHERE ii.invoice_id = ?
	`

	rows, err := db.Query(itemQuery, id)
	if err != nil {
		fmt.Println("err", err)
		return inv, nil, err
	}
	defer rows.Close()

	var items []InvoiceItem

	for rows.Next() {
		var item InvoiceItem

		err := rows.Scan(
			&item.ID,
			&item.InvoiceID,
			&item.ProductID,
			&item.ProductName,
			&item.Qty,
			&item.Price,
			&item.GST,
			&item.LineTotal,
		)

		if err != nil {
			return inv, nil, err
		}

		items = append(items, item)
	}

	// 3. Convert to view model AFTER loop
	var itemsView []InvoiceItemView

	for i, it := range items {
		itemsView = append(itemsView, InvoiceItemView{
			No:          i + 1,
			InvoiceItem: it,
		})
	}

	return inv, itemsView, nil
}

func (app *App) InvoiceHandler(w http.ResponseWriter, r *http.Request) {

	invoices, _ := getInvoices()

	data := map[string]interface{}{
		"Title":    "Invoice",
		"Page":     "invoice",
		"Invoices": invoices,
	}

	app.Tmpl.ExecuteTemplate(w, "invoice", data)
}

func (app *App) CreateInvoiceHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	company, err := getCompany()
	if err != nil {
		http.Error(w, "company not found", 500)
		return
	}

	// SHOW PAGE
	if r.Method == http.MethodGet {

		customers, _ := getCustomers()

		products, _ := getProducts()

		data := map[string]interface{}{
			"Title":     "Create Invoice",
			"Page":      "create_invoice",
			"Customers": customers,
			"Products":  products,
		}

		app.Tmpl.ExecuteTemplate(w, "invoice", data)
		return
	}

	// SAVE INVOICE
	if r.Method == http.MethodPost {

		r.ParseForm()

		customerID := r.FormValue("customer_id")
		productIDs := r.Form["product_id[]"]
		qtys := r.Form["qty[]"]

		if len(productIDs) != len(qtys) {
			http.Error(w, "invalid form data", 400)
			return
		}

		// start transaction (IMPORTANT)
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// totals
		var subtotal, gstTotal float64

		var invoiceNo string = generateInvoiceNumber(company.CompanyName)

		// create invoice first
		res, err := tx.Exec(
			`INSERT INTO invoices(invoice_no, customer_id, subtotal, gst_total, grand_total, created_at)
		 VALUES(?, ?,?,?,?,datetime('now'))`,
			invoiceNo, customerID, 0, 0, 0,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), 500)
			return
		}

		invoiceID, _ := res.LastInsertId()

		// loop items
		for i := range productIDs {

			productID := productIDs[i]
			qty := qtys[i]

			var price, gst float64

			// get product details
			err := db.QueryRow(
				"SELECT price, gst FROM products WHERE id=?",
				productID,
			).Scan(&price, &gst)

			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), 500)
				return
			}

			qtyInt := atoi(qty)

			lineTotal := price * float64(qtyInt)
			gstAmt := (lineTotal * gst) / 100

			subtotal += lineTotal
			gstTotal += gstAmt

			_, err = tx.Exec(
				`INSERT INTO invoice_items(invoice_id, product_id, qty, price, gst, line_total)
			 VALUES(?,?,?,?,?,?)`,
				invoiceID, productID, qtyInt, price, gst, lineTotal+gstAmt,
			)

			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), 500)
				return
			}
		}

		grandTotal := subtotal + gstTotal

		// update invoice totals
		_, err = tx.Exec(
			`UPDATE invoices 
		 SET subtotal=?, gst_total=?, grand_total=? 
		 WHERE id=?`,
			subtotal, gstTotal, grandTotal, invoiceID,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), 500)
			return
		}

		tx.Commit()

		http.Redirect(w, r, "/invoice", http.StatusSeeOther)
	}
}

func (app *App) GetInvoiceViewHandler(w http.ResponseWriter, r *http.Request) {

	// allow only GET
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get query param: /invoice/view?id=2
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// call DB function
	inv, items, err := getInvoiceByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Invoice",
		"Page":    "view_invoice",
		"Invoice": inv,
		"Items":   items,
	}

	app.Tmpl.ExecuteTemplate(w, "invoice", data)
}
