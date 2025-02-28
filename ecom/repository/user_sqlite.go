package repository

import (
	"database/sql"

	"ecom.com/models"
)

type SQLiteUserRepository struct {
	DB *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{DB: db}
}

func (r *SQLiteUserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, email) VALUES (?, ?)`
	_, err := r.DB.Exec(query, user.Name, user.Email)
	return err
}

func (r *SQLiteUserRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, name, email FROM users`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
