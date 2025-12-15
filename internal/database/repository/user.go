package repository

import (
	"database/sql"
	"trmp/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users WHERE email = ?
	`

	row := r.db.QueryRow(query, email)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*model.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
