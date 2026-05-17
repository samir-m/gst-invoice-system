package main

type Customer struct {
	ID    int
	Name  string
	Phone string
	GSTIN string
}

type Product struct {
	ID      int
	Product string
	Price   float64
	GST     float64
}

type Invoice struct {
	ID         int
	InvoiceNo  string
	CustomerID int
	Subtotal   float64
	GSTTotal   float64
	GrandTotal float64
	Date       string
}

type InvoiceItem struct {
	ID          int
	InvoiceID   int
	ProductID   int
	ProductName string // optional (denormalized for display)
	Qty         int
	Price       float64
	GST         float64
	LineTotal   float64
}

type InvoiceList struct {
	ID         int
	InvoiceNo  string
	Customer   string
	Subtotal   float64
	GSTTotal   float64
	GrandTotal float64
	Date       string
}
