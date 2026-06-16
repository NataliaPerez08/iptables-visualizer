package deployment

import (
	"fmt"
	"os/exec"
	"strings"
)

type DeployResult struct {
	Success  bool     `json:"success"`
	Commands []string `json:"commands"`
	Output   string   `json:"output"`
	Error    string   `json:"error,omitempty"`
}

type Deployer struct {
	DryRunOnly bool
}

func NewDeployer(dryRunOnly bool) *Deployer {
	return &Deployer{DryRunOnly: dryRunOnly}
}

func (d *Deployer) Apply(rules []string) *DeployResult {
	result := &DeployResult{Commands: rules}

	if d.DryRunOnly {
		result.Success = true
		result.Output = "DRY-RUN: commands not executed (dry-run mode enabled)"
		return result
	}

	var outputs []string
	for _, cmd := range rules {
		out, err := d.execCommand(cmd)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("command %q failed: %s", cmd, err.Error())
			result.Output = strings.Join(outputs, "\n")
			return result
		}
		outputs = append(outputs, out)
	}

	result.Success = true
	result.Output = strings.Join(outputs, "\n")
	return result
}

func (d *Deployer) Rollback(rules []string) *DeployResult {
	var rollbackCmds []string
	for _, cmd := range rules {
		rollback := d.invertCommand(cmd)
		if rollback != "" {
			rollbackCmds = append(rollbackCmds, rollback)
		}
	}
	return d.Apply(rollbackCmds)
}

func (d *Deployer) execCommand(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	name := parts[0]
	args := parts[1:]
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s: %w", string(out), err)
	}
	return string(out), nil
}

func (d *Deployer) invertCommand(cmd string) string {
	if strings.HasPrefix(cmd, "iptables -A") {
		return strings.Replace(cmd, "iptables -A", "iptables -D", 1)
	}
	if strings.HasPrefix(cmd, "nft add rule") {
		return strings.Replace(cmd, "nft add rule", "nft delete rule", 1)
	}
	return ""
}
