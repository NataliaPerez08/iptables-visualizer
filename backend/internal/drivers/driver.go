package drivers

import "github.com/anomalyco/iptables-visualizer/internal/models"

type Driver interface {
	Name() string
	CompileRule(rule *models.Rule) (string, error)
}
