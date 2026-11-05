package repository

import (
	"expenses/internal/models"
	"fmt"
)

func (repo ExpenseRepository) GetCategoryStats(userID int, from string, to string) ([]models.CategoryStat, error) {
	stats := []models.CategoryStat{}
	query := "SELECT categories.name, SUM(expenses.amount) FROM expenses JOIN categories ON expenses.category_id = categories.id WHERE expenses.user_id = ?"
	args := []any{userID}
	if from != "" {
		query += " AND expenses.date >= ?"
		args = append(args, from)
	}
	if to != "" {
		query += " AND expenses.date <= ?"
		args = append(args, to)
	}
	query += " GROUP BY categories.id, categories.name ORDER BY categories.id"
	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var stat models.CategoryStat
		err = rows.Scan(&stat.Name, &stat.Total)
		if err != nil {
			return stats, err
		}
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}

func (repo ExpenseRepository) GetMonthStats(userID int, year int) ([]models.MonthStat, error) {
	stats := []models.MonthStat{}
	rows, err := repo.DB.Query("SELECT CAST(strftime('%m', date) AS INTEGER), SUM(amount) FROM expenses WHERE user_id = ? AND strftime('%Y', date) = ? GROUP BY strftime('%m', date) ORDER BY strftime('%m', date)", userID, fmt.Sprintf("%04d", year))
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var stat models.MonthStat
		err = rows.Scan(&stat.Month, &stat.Total)
		if err != nil {
			return stats, err
		}
		stats = append(stats, stat)
	}
	return stats, rows.Err()
}
