package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type UserRepository interface {
	FindByID(id int64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	List(limit, offset int) ([]models.User, error)
	Delete(id int64) error
}

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) FindByID(id int64) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, COALESCE(email,''), is_active, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

func (r *SQLiteUserRepository) FindByUsername(username string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, role, COALESCE(email,''), is_active, created_at, updated_at
		 FROM users WHERE username = ?`, username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	return u, err
}

func (r *SQLiteUserRepository) Create(u *models.User) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	res, err := r.db.Exec(
		`INSERT INTO users (username, password_hash, role, email, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.Username, u.PasswordHash, u.Role, u.Email, u.IsActive, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.ID = id
	return nil
}

func (r *SQLiteUserRepository) Update(u *models.User) error {
	u.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		`UPDATE users SET username=?, role=?, email=?, is_active=?, updated_at=? WHERE id=?`,
		u.Username, u.Role, u.Email, u.IsActive, u.UpdatedAt, u.ID,
	)
	return err
}

func (r *SQLiteUserRepository) List(limit, offset int) ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, password_hash, role, COALESCE(email,''), is_active, created_at, updated_at
		 FROM users ORDER BY id LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *SQLiteUserRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
