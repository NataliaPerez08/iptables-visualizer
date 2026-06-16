package models

import "time"

type PolicyStatus string

const (
	PolicyDraft     PolicyStatus = "draft"
	PolicyActive    PolicyStatus = "active"
	PolicyInactive  PolicyStatus = "inactive"
	PolicyFailed    PolicyStatus = "failed"
)

type MatchProtocol string

const (
	ProtocolTCP  MatchProtocol = "tcp"
	ProtocolUDP  MatchProtocol = "udp"
	ProtocolICMP MatchProtocol = "icmp"
	ProtocolAny  MatchProtocol = "any"
)

type RuleAction string

const (
	ActionAccept RuleAction = "accept"
	ActionDrop   RuleAction = "drop"
	ActionReject RuleAction = "reject"
	ActionLog    RuleAction = "log"
)

type Rule struct {
	ID          string        `json:"id"`
	Order       int           `json:"order"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Enabled     bool          `json:"enabled"`
	Action      RuleAction    `json:"action"`
	Protocol    MatchProtocol `json:"protocol"`
	SrcAddr     string        `json:"src_addr"`
	SrcPort     string        `json:"src_port,omitempty"`
	DstAddr     string        `json:"dst_addr"`
	DstPort     string        `json:"dst_port,omitempty"`
	InInterface string        `json:"in_interface,omitempty"`
	OutInterface string       `json:"out_interface,omitempty"`
	State       string        `json:"state,omitempty"`
	LogPrefix   string        `json:"log_prefix,omitempty"`
	Extra       []string      `json:"extra,omitempty"`
}

type Policy struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Rules       []Rule        `json:"rules"`
	Status      PolicyStatus  `json:"status"`
	Version     int           `json:"version"`
	CreatedBy   int64         `json:"created_by"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	AppliedAt   *time.Time    `json:"applied_at,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
}

type CreatePolicyRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Rules       []Rule   `json:"rules"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdatePolicyRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Rules       []Rule   `json:"rules,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
