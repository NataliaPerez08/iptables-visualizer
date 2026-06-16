package engine

import (
	"net"
	"regexp"
	"strings"

	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Issues []ValidationIssue `json:"issues"`
}

type ValidationIssue struct {
	RuleID  string `json:"rule_id,omitempty"`
	Field   string `json:"field"`
	Message string `json:"message"`
	Severity string `json:"severity"`
}

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

var validActions = map[models.RuleAction]bool{
	models.ActionAccept: true,
	models.ActionDrop:   true,
	models.ActionReject: true,
	models.ActionLog:    true,
}

var validProtocols = map[models.MatchProtocol]bool{
	models.ProtocolTCP:  true,
	models.ProtocolUDP:  true,
	models.ProtocolICMP: true,
	models.ProtocolAny:  true,
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)
var portRegex = regexp.MustCompile(`^(\d{1,5})(:\d{1,5})?$`)

func (v *Validator) Validate(policy *models.Policy) ValidationResult {
	result := ValidationResult{Valid: true}

	if strings.TrimSpace(policy.Name) == "" {
		result.Issues = append(result.Issues, ValidationIssue{
			Field:    "name",
			Message:  "policy name is required",
			Severity: "error",
		})
	}

	if !nameRegex.MatchString(policy.Name) {
		result.Issues = append(result.Issues, ValidationIssue{
			Field:    "name",
			Message:  "policy name must be 1-128 chars, alphanumeric with _ and -",
			Severity: "error",
		})
	}

	if len(policy.Rules) == 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Field:    "rules",
			Message:  "policy must contain at least one rule",
			Severity: "warning",
		})
	}

	for i := range policy.Rules {
		rule := &policy.Rules[i]
		v.validateRule(rule, &result)
	}

	result.Valid = len(result.Issues) == 0
	for _, iss := range result.Issues {
		if iss.Severity == "error" {
			result.Valid = false
			break
		}
	}

	return result
}

func (v *Validator) validateRule(rule *models.Rule, result *ValidationResult) {
	if strings.TrimSpace(rule.Name) == "" {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "name",
			Message:  "rule name is required",
			Severity: "error",
		})
	}

	if !validActions[rule.Action] {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "action",
			Message:  "invalid action: must be accept, drop, reject, or log",
			Severity: "error",
		})
	}

	if !validProtocols[rule.Protocol] {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "protocol",
			Message:  "invalid protocol: must be tcp, udp, icmp, or any",
			Severity: "error",
		})
	}

	if rule.SrcAddr != "" && !isValidAddr(rule.SrcAddr) {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "src_addr",
			Message:  "invalid source address",
			Severity: "error",
		})
	}

	if rule.DstAddr != "" && !isValidAddr(rule.DstAddr) {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "dst_addr",
			Message:  "invalid destination address",
			Severity: "error",
		})
	}

	if rule.SrcPort != "" && !isValidPort(rule.SrcPort) {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "src_port",
			Message:  "invalid source port",
			Severity: "error",
		})
	}

	if rule.DstPort != "" && !isValidPort(rule.DstPort) {
		result.Issues = append(result.Issues, ValidationIssue{
			RuleID:   rule.ID,
			Field:    "dst_port",
			Message:  "invalid destination port",
			Severity: "error",
		})
	}
}

func isValidAddr(addr string) bool {
	if strings.HasPrefix(addr, "! ") {
		addr = addr[2:]
	}
	if strings.Contains(addr, "/") {
		_, _, err := net.ParseCIDR(addr)
		return err == nil
	}
	return net.ParseIP(addr) != nil || addr == "0.0.0.0/0"
}

func isValidPort(port string) bool {
	if strings.HasPrefix(port, "! ") {
		port = port[2:]
	}
	return portRegex.MatchString(port)
}
