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
	return err
}

func (repo ExpenseRepository) Create(expense models.Expense) error {
	_, err := repo.DB.Exec("INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)", expense.Amount, expense.Description, expense.Date)
	return err
}

func (repo ExpenseRepository) GetAll() ([]models.Expense, error) {
	expenses := []models.Expense{}
	rows, err := repo.DB.Query("SELECT id, amount, description, date FROM expenses ORDER BY id")
	if err != nil {
		return expenses, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense
		err = rows.Scan(&expense.ID, &expense.Amount, &expense.Description, &expense.Date)
		if err != nil {
			return expenses, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, rows.Err()
}

func (repo ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var expense models.Expense
	err := repo.DB.QueryRow("SELECT id, amount, description, date FROM expenses WHERE id = ?", id).Scan(&expense.ID, &expense.Amount, &expense.Description, &expense.Date)
	return expense, err
}

func (repo ExpenseRepository) Update(expense models.Expense) error {
	_, err := repo.DB.Exec("UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?", expense.Amount, expense.Description, expense.Date, expense.ID)
	return err
}

func (repo ExpenseRepository) Delete(id int) error {
	_, err := repo.DB.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}
