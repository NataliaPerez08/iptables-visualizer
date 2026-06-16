package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type PolicyRepository interface {
	FindByID(id int64) (*models.Policy, error)
	Create(policy *models.Policy) error
	Update(policy *models.Policy) error
	Delete(id int64) error
	List(limit, offset int) ([]models.Policy, error)
	UpdateStatus(id int64, status models.PolicyStatus) error
}

type SQLitePolicyRepository struct {
	db *sql.DB
}

func NewPolicyRepository(db *sql.DB) PolicyRepository {
	return &SQLitePolicyRepository{db: db}
}

func (r *SQLitePolicyRepository) FindByID(id int64) (*models.Policy, error) {
	p := &models.Policy{}
	var rulesJSON string
	var tagsJSON string
	var appliedAt sql.NullTime

	err := r.db.QueryRow(
		`SELECT id, name, COALESCE(description,''), rules, status, version, created_by, created_at, updated_at, applied_at, COALESCE(tags,'[]')
		 FROM policies WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &rulesJSON, &p.Status, &p.Version, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt, &appliedAt, &tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("policy not found: %w", err)
	}

	if err := json.Unmarshal([]byte(rulesJSON), &p.Rules); err != nil {
		return nil, fmt.Errorf("failed to parse rules: %w", err)
	}
	if err := json.Unmarshal([]byte(tagsJSON), &p.Tags); err != nil {
		p.Tags = nil
	}
	if appliedAt.Valid {
		p.AppliedAt = &appliedAt.Time
	}
	return p, nil
}

func (r *SQLitePolicyRepository) Create(p *models.Policy) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Version = 1

	rulesJSON, _ := json.Marshal(p.Rules)
	tagsJSON, _ := json.Marshal(p.Tags)

	res, err := r.db.Exec(
		`INSERT INTO policies (name, description, rules, status, version, created_by, created_at, updated_at, tags)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.Description, string(rulesJSON), p.Status, p.Version, p.CreatedBy, p.CreatedAt, p.UpdatedAt, string(tagsJSON),
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	return nil
}

func (r *SQLitePolicyRepository) Update(p *models.Policy) error {
	p.UpdatedAt = time.Now()
	p.Version++

	rulesJSON, _ := json.Marshal(p.Rules)
	tagsJSON, _ := json.Marshal(p.Tags)

	_, err := r.db.Exec(
		`UPDATE policies SET name=?, description=?, rules=?, status=?, version=?, updated_at=?, tags=? WHERE id=?`,
		p.Name, p.Description, string(rulesJSON), p.Status, p.Version, p.UpdatedAt, string(tagsJSON), p.ID,
	)
	return err
}

func (r *SQLitePolicyRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM policies WHERE id = ?`, id)
	return err
}

func (r *SQLitePolicyRepository) List(limit, offset int) ([]models.Policy, error) {
	rows, err := r.db.Query(
		`SELECT id, name, COALESCE(description,''), rules, status, version, created_by, created_at, updated_at, applied_at, COALESCE(tags,'[]')
		 FROM policies ORDER BY updated_at DESC LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []models.Policy
	for rows.Next() {
		var p models.Policy
		var rulesJSON string
		var tagsJSON string
		var appliedAt sql.NullTime

		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &rulesJSON, &p.Status, &p.Version, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt, &appliedAt, &tagsJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(rulesJSON), &p.Rules); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsJSON), &p.Tags)
		if appliedAt.Valid {
			p.AppliedAt = &appliedAt.Time
		}
		policies = append(policies, p)
	}
	return policies, rows.Err()
}

func (r *SQLitePolicyRepository) UpdateStatus(id int64, status models.PolicyStatus) error {
	now := time.Now()
	var appliedAt interface{}
	if status == models.PolicyActive {
		appliedAt = now
	}
	_, err := r.db.Exec(
		`UPDATE policies SET status=?, applied_at=?, updated_at=? WHERE id=?`,
		status, appliedAt, now, id,
	)
	return err
}
