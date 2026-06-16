package drivers

import (
	"fmt"
	"strings"

	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type NFTablesDriver struct{}

func (d *NFTablesDriver) Name() string { return "nftables" }

func (d *NFTablesDriver) CompileRule(rule *models.Rule) (string, error) {
	parts := []string{"nft", "add", "rule", "inet", "filter", "INPUT"}

	switch rule.Action {
	case models.ActionAccept:
		parts = append(parts, "accept")
	case models.ActionDrop:
		parts = append(parts, "drop")
	case models.ActionReject:
		parts = append(parts, "reject")
	case models.ActionLog:
		parts = append(parts, "log")
	default:
		return "", fmt.Errorf("unsupported action: %s", rule.Action)
	}

	if rule.Protocol != "" && rule.Protocol != models.ProtocolAny {
		meta := fmt.Sprintf("meta l4proto %s", rule.Protocol)
		parts = append(parts, meta)
	}

	if rule.SrcAddr != "" {
		parts = append(parts, "ip", "saddr", rule.SrcAddr)
	}

	if rule.DstAddr != "" {
		parts = append(parts, "ip", "daddr", rule.DstAddr)
	}

	if rule.SrcPort != "" || rule.DstPort != "" {
		if rule.Protocol == models.ProtocolTCP || rule.Protocol == models.ProtocolUDP {
			if rule.SrcPort != "" {
				parts = append(parts, "th", "sport", rule.SrcPort)
			}
			if rule.DstPort != "" {
				parts = append(parts, "th", "dport", rule.DstPort)
			}
		}
	}

	if rule.InInterface != "" {
		parts = append(parts, "iif", rule.InInterface)
	}

	if rule.OutInterface != "" {
		parts = append(parts, "oif", rule.OutInterface)
	}

	if rule.State != "" {
		parts = append(parts, "ct", "state", strings.ToLower(rule.State))
	}

	if rule.LogPrefix != "" {
		parts = append(parts, "log", "prefix", fmt.Sprintf("%q", rule.LogPrefix))
	}

	if len(rule.Extra) > 0 {
		parts = append(parts, rule.Extra...)
	}

	return strings.Join(parts, " "), nil
}
