package repository

import (
	"database/sql"
	"time"

	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type AuditRepository interface {
	Create(log *models.AuditLog) error
	Query(query models.AuditQuery) ([]models.AuditLog, error)
}

type SQLiteAuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) AuditRepository {
	return &SQLiteAuditRepository{db: db}
}

func (r *SQLiteAuditRepository) Create(l *models.AuditLog) error {
	l.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO audit_logs (user_id, username, action, resource, details, ip_address, success, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		l.UserID, l.Username, l.Action, l.Resource, l.Details, l.IPAddress, l.Success, l.CreatedAt,
	)
	return err
}

func (r *SQLiteAuditRepository) Query(q models.AuditQuery) ([]models.AuditLog, error) {
	query := `SELECT id, user_id, username, action, resource, COALESCE(details,''), COALESCE(ip_address,''), success, created_at
			  FROM audit_logs WHERE 1=1`
	args := []interface{}{}

	if q.UserID > 0 {
		query += ` AND user_id = ?`
		args = append(args, q.UserID)
	}
	if q.Action != "" {
		query += ` AND action = ?`
		args = append(args, q.Action)
	}
	if q.Resource != "" {
		query += ` AND resource = ?`
		args = append(args, q.Resource)
	}
	if q.From != nil {
		query += ` AND created_at >= ?`
		args = append(args, *q.From)
	}
	if q.To != nil {
		query += ` AND created_at <= ?`
		args = append(args, *q.To)
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	limit := q.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := q.Offset
	if offset < 0 {
		offset = 0
	}
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Username, &l.Action, &l.Resource, &l.Details, &l.IPAddress, &l.Success, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
