package deployment

import (
	"fmt"
	"strings"

	"github.com/anomalyco/iptables-visualizer/internal/models"
	"github.com/anomalyco/iptables-visualizer/internal/engine"
)

type DryRunReport struct {
	PolicyName      string            `json:"policy_name"`
	PolicyID        int64             `json:"policy_id"`
	Driver          string            `json:"driver"`
	RuleCount       int               `json:"rule_count"`
	EnabledRules    int               `json:"enabled_rules"`
	Commands        []string          `json:"commands"`
	Validation      engine.ValidationResult `json:"validation"`
	Changes         *PolicyDiff       `json:"changes,omitempty"`
}

type PolicyDiff struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Unchanged int    `json:"unchanged"`
}

func GenerateDryRun(policy *models.Policy, driverName string, compiler *engine.Compiler, validator *engine.Validator, previousCommands []string) *DryRunReport {
	report := &DryRunReport{
		PolicyName:   policy.Name,
		PolicyID:     policy.ID,
		Driver:       driverName,
		RuleCount:    len(policy.Rules),
		EnabledRules: countEnabled(policy.Rules),
	}

	report.Validation = validator.Validate(policy)

	commands, err := compiler.Compile(policy, driverName)
	if err != nil {
		report.Commands = []string{fmt.Sprintf("ERROR: %s", err.Error())}
		return report
	}
	report.Commands = commands

	if len(previousCommands) > 0 {
		report.Changes = diffCommands(previousCommands, commands)
	}

	return report
}

func countEnabled(rules []models.Rule) int {
	count := 0
	for _, r := range rules {
		if r.Enabled {
			count++
		}
	}
	return count
}

func diffCommands(old, new []string) *PolicyDiff {
	oldSet := make(map[string]bool)
	for _, c := range old {
		oldSet[c] = true
	}
	newSet := make(map[string]bool)
	for _, c := range new {
		newSet[c] = true
	}

	diff := &PolicyDiff{}
	for _, c := range new {
		if !oldSet[c] {
			diff.Added = append(diff.Added, c)
		}
	}
	for _, c := range old {
		if !newSet[c] {
			diff.Removed = append(diff.Removed, c)
		}
	}
	diff.Unchanged = len(old) - len(diff.Removed)
	if diff.Unchanged < 0 {
		diff.Unchanged = 0
	}
	return diff
}

func FormatDryRunReport(report *DryRunReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Policy: %s (ID: %d)\n", report.PolicyName, report.PolicyID))
	sb.WriteString(fmt.Sprintf("Driver: %s\n", report.Driver))
	sb.WriteString(fmt.Sprintf("Total rules: %d (%d enabled)\n", report.RuleCount, report.EnabledRules))
	sb.WriteString(fmt.Sprintf("Validation: %t\n\n", report.Validation.Valid))

	if !report.Validation.Valid {
		sb.WriteString("Validation Issues:\n")
		for _, iss := range report.Validation.Issues {
			sb.WriteString(fmt.Sprintf("  - [%s] %s: %s\n", iss.Severity, iss.Field, iss.Message))
		}
		sb.WriteString("\n")
	}

	if report.Changes != nil {
		sb.WriteString("Changes vs previous:\n")
		sb.WriteString(fmt.Sprintf("  Added: %d\n", len(report.Changes.Added)))
		sb.WriteString(fmt.Sprintf("  Removed: %d\n", len(report.Changes.Removed)))
		sb.WriteString(fmt.Sprintf("  Unchanged: %d\n\n", report.Changes.Unchanged))
	}

	sb.WriteString("Commands to execute:\n")
	for _, cmd := range report.Commands {
		sb.WriteString(fmt.Sprintf("  %s\n", cmd))
	}

	return sb.String()
}
