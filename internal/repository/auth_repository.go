package repository

import (
	"database/sql"
	"task-management/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"errors"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create (user models.User) error {
	query := `
		INSERT INTO users
		(name, email, password)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, user.Name, user.Email, user.Password)

	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique constraint violation
			return errors.New("email already exists")
		}
		return err
	}

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