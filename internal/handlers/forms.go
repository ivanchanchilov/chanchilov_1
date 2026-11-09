package handlers

import (
	"expenses/internal/models"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ReadExpense(r *http.Request, categories []models.Category) (models.Expense, string) {
	var expense models.Expense
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 {
		return expense, "Сумма должна быть положительным числом."
	}
	expense.Amount = amount
	expense.Description = strings.TrimSpace(r.FormValue("description"))
	if expense.Description == "" {
		return expense, "Введите описание траты."
	}
	expense.Date = r.FormValue("date")
	_, err = time.Parse("2006-01-02", expense.Date)
	if err != nil {
		return expense, "Введите дату в формате ГГГГ-ММ-ДД."
	}
	expense.CategoryID, err = strconv.Atoi(r.FormValue("category"))
	if err != nil {
		return expense, "Выберите категорию из списка."
	}
	found := false
	for _, category := range categories {
		if category.ID == expense.CategoryID {
			found = true
		}
	}
	if !found {
		return expense, "Выберите категорию из списка."
	}
	return expense, ""
}

func ExpenseForm(w http.ResponseWriter, r *http.Request, name string, expense models.Expense, categories []models.Category, message string) {
	data := models.PageData{
		Expense:         expense,
		Categories:      categories,
		Error:           message,
		FormAmount:      strconv.FormatFloat(expense.Amount, 'f', -1, 64),
		FormDescription: expense.Description,
		FormDate:        expense.Date,
		FormCategory:    strconv.Itoa(expense.CategoryID),
	}
	if r.Method == "POST" {
		data.FormAmount = r.FormValue("amount")
		data.FormDescription = r.FormValue("description")
		data.FormDate = r.FormValue("date")
		data.FormCategory = r.FormValue("category")
	}
	found := false
	for _, category := range categories {
		if strconv.Itoa(category.ID) == data.FormCategory {
			found = true
		}
	}
	data.InvalidCategory = !found
	if r.Method == "GET" && name == "add" {
		data.FormAmount = ""
		data.FormCategory = "1"
		data.InvalidCategory = false
	}
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/"+name+".html")
	if err != nil {
		ServerError(w)
		return
	}
	if message != "" {
		w.WriteHeader(400)
	}
	t.ExecuteTemplate(w, "layout", data)
}
