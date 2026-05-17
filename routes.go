package main

import "net/http"

func NewRouter(app *App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/customer", app.CustomerHandler)
	mux.HandleFunc("/customer/add", app.AddCustomer)

	mux.HandleFunc("/product", app.ProductHandler)
	mux.HandleFunc("/product/add", app.AddProductHandler)

	mux.HandleFunc("/invoice", app.InvoiceHandler)
	mux.HandleFunc("/invoice/create", app.CreateInvoiceHandler)

	mux.HandleFunc("/invoice/view", app.GetInvoiceViewHandler)

	mux.HandleFunc("/invoice/pdf", app.ShowInvoicePdf)

	mux.HandleFunc("/invoice/print", app.InvoicePDFHandler)

	return mux
}
