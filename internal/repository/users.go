package repository

import "expenses/internal/models"

func (repo ExpenseRepository) CreateUser(email string, passwordHash string) error {
	_, err := repo.DB.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, passwordHash)
	return err
}

func (repo ExpenseRepository) GetUser(email string) (models.User, error) {
	var user models.User
	err := repo.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	return user, err
}

func (repo ExpenseRepository) CreateSession(sessionID string, userID int, expiresAt int64) error {
	_, err := repo.DB.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	return err
}

func (repo ExpenseRepository) GetSessionUser(sessionID string, now int64) (int, error) {
	var userID int
	err := repo.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ? AND expires_at > ?", sessionID, now).Scan(&userID)
	return userID, err
}

func (repo ExpenseRepository) DeleteSession(sessionID string) error {
	_, err := repo.DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}
