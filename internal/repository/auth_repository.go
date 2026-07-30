package repository

import (
	"database/sql"
	"fmt"
	"task-management/internal/models"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}


func (r *AuthRepository) EmailExists(email string) (bool, error) {
	query := `
		SELECT 1 FROM users WHERE email = $1
	`

	err := r.db.QueryRow(query, email)
	if err != nil {
		fmt.Println("wweewfnjkn")
		return true, err.Err()
	}

	return false, nil
}


func (r *AuthRepository) Create (user models.User) error {
	query := `
		INSERT INTO users
		(name, email, password)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, user.Name, user.Email, user.Password)

	return err
}

func (r *AuthRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT 
			id, name, email, password, created_at, updated_at
		From users
		WHERE email = $1
	`
	var user models.User

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	
	return &user, nil
}