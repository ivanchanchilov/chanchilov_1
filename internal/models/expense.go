package models

type Expense struct {
	UserID       int
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

type User struct {
	ID           int
	Email        string
	PasswordHash string
}
