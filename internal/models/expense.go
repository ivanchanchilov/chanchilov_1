package models

type Expense struct {
	ID           int
	CategoryID   int
	CategoryName string
	Amount       float64
	Description  string
	Date         string
}

type Category struct {
	ID   int
	Name string
}

type PageData struct {
	Expenses   []Expense
	Expense    Expense
	Categories []Category
	CategoryID int
	From       string
	To         string
}
