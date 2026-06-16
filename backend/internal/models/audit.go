package models

import "time"

type AuditAction string

const (
	AuditLogin        AuditAction = "login"
	AuditLogout       AuditAction = "logout"
	AuditCreatePolicy AuditAction = "create_policy"
	AuditUpdatePolicy AuditAction = "update_policy"
	AuditDeletePolicy AuditAction = "delete_policy"
	AuditApplyPolicy  AuditAction = "apply_policy"
	AuditDryRunPolicy AuditAction = "dry_run_policy"
	AuditValidatePolicy AuditAction = "validate_policy"
	AuditCreateUser   AuditAction = "create_user"
	AuditUpdateUser   AuditAction = "update_user"
	AuditDeleteUser   AuditAction = "delete_user"
)

type AuditLog struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	Username  string      `json:"username"`
	Action    AuditAction `json:"action"`
	Resource  string      `json:"resource"`
	Details   string      `json:"details,omitempty"`
	IPAddress string      `json:"ip_address,omitempty"`
	Success   bool        `json:"success"`
	CreatedAt time.Time   `json:"created_at"`
}

type AuditQuery struct {
	UserID   int64      `json:"user_id,omitempty"`
	Action   AuditAction `json:"action,omitempty"`
	Resource string     `json:"resource,omitempty"`
	From     *time.Time `json:"from,omitempty"`
	To       *time.Time `json:"to,omitempty"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}
