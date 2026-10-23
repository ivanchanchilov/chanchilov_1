package repository

import (
	"database/sql"
	"expenses/internal/models"
	"os"
)

type ExpenseRepository struct {
	DB *sql.DB
}

func (repo ExpenseRepository) Init() error {
	migration, err := os.ReadFile("migrations/001_create_expenses.sql")
	if err != nil {
		return err
	}
	_, err = repo.DB.Exec(string(migration))
	if err != nil {
		return err
	}
	var version int
	err = repo.DB.QueryRow("PRAGMA user_version").Scan(&version)
	if err != nil {
		return err
	}
	if version < 2 {
		migration, err = os.ReadFile("migrations/002_add_categories.sql")
		if err != nil {
			return err
		}
		_, err = repo.DB.Exec(string(migration))
	}
	return err
}

func (repo ExpenseRepository) Create(expense models.Expense) error {
	_, err := repo.DB.Exec("INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)", expense.Amount, expense.Description, expense.Date, expense.CategoryID)
	return err
}

func (repo ExpenseRepository) GetAll(categoryID int, from string, to string) ([]models.Expense, error) {
	expenses := []models.Expense{}
	query := "SELECT expenses.id, expenses.amount, expenses.description, expenses.date, expenses.category_id, categories.name FROM expenses JOIN categories ON expenses.category_id = categories.id WHERE 1 = 1"
	args := []any{}
	if categoryID != 0 {
		query += " AND expenses.category_id = ?"
		args = append(args, categoryID)
	}
	if from != "" {
		query += " AND expenses.date >= ?"
		args = append(args, from)
	}
	if to != "" {
		query += " AND expenses.date <= ?"
		args = append(args, to)
	}
	query += " ORDER BY expenses.id"
	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return expenses, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense
		err = rows.Scan(&expense.ID, &expense.Amount, &expense.Description, &expense.Date, &expense.CategoryID, &expense.CategoryName)
		if err != nil {
			return expenses, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, rows.Err()
}

func (repo ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var expense models.Expense
	err := repo.DB.QueryRow("SELECT expenses.id, expenses.amount, expenses.description, expenses.date, expenses.category_id, categories.name FROM expenses JOIN categories ON expenses.category_id = categories.id WHERE expenses.id = ?", id).Scan(&expense.ID, &expense.Amount, &expense.Description, &expense.Date, &expense.CategoryID, &expense.CategoryName)
	return expense, err
}

func (repo ExpenseRepository) Update(expense models.Expense) error {
	_, err := repo.DB.Exec("UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?", expense.Amount, expense.Description, expense.Date, expense.CategoryID, expense.ID)
	return err
}

func (repo ExpenseRepository) Delete(id int) error {
	_, err := repo.DB.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}

func (repo ExpenseRepository) GetCategories() ([]models.Category, error) {
	categories := []models.Category{}
	rows, err := repo.DB.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
