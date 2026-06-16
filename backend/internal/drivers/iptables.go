package drivers

import (
	"fmt"
	"strings"

	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type IPTablesDriver struct{}

func (d *IPTablesDriver) Name() string { return "iptables" }

func (d *IPTablesDriver) CompileRule(rule *models.Rule) (string, error) {
	parts := []string{"iptables"}

	switch rule.Action {
	case models.ActionAccept:
		parts = append(parts, "-A", "INPUT", "-j", "ACCEPT")
	case models.ActionDrop:
		parts = append(parts, "-A", "INPUT", "-j", "DROP")
	case models.ActionReject:
		parts = append(parts, "-A", "INPUT", "-j", "REJECT")
	case models.ActionLog:
		parts = append(parts, "-A", "INPUT", "-j", "LOG")
	default:
		return "", fmt.Errorf("unsupported action: %s", rule.Action)
	}

	if rule.Protocol != "" && rule.Protocol != models.ProtocolAny {
		parts = append(parts, "-p", string(rule.Protocol))
	}

	if rule.SrcAddr != "" {
		parts = append(parts, "-s", rule.SrcAddr)
	}

	if rule.DstAddr != "" {
		parts = append(parts, "-d", rule.DstAddr)
	}

	if rule.SrcPort != "" {
		parts = append(parts, "--sport", rule.SrcPort)
	}

	if rule.DstPort != "" {
		parts = append(parts, "--dport", rule.DstPort)
	}

	if rule.InInterface != "" {
		parts = append(parts, "-i", rule.InInterface)
	}

	if rule.OutInterface != "" {
		parts = append(parts, "-o", rule.OutInterface)
	}

	if rule.State != "" {
		parts = append(parts, "-m", "state", "--state", rule.State)
	}

	if rule.LogPrefix != "" {
		parts = append(parts, "--log-prefix", fmt.Sprintf("%q", rule.LogPrefix))
	}

	if len(rule.Extra) > 0 {
		parts = append(parts, rule.Extra...)
	}

	return strings.Join(parts, " "), nil
}
