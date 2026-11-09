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
	InvalidCategory bool
	Error           string
	FormAmount      string
	FormDescription string
	FormDate        string
	FormCategory    string
	Expenses        []Expense
	Expense         Expense
	Categories      []Category
	CategoryID      int
	From            string
	To              string
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
}
