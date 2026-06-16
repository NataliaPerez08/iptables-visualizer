package engine

import (
	"fmt"

	"github.com/anomalyco/iptables-visualizer/internal/drivers"
	"github.com/anomalyco/iptables-visualizer/internal/models"
)

type Compiler struct {
	drivers map[string]drivers.Driver
}

func NewCompiler() *Compiler {
	return &Compiler{
		drivers: map[string]drivers.Driver{
			"iptables": &drivers.IPTablesDriver{},
			"nftables": &drivers.NFTablesDriver{},
		},
	}
}

func (c *Compiler) Compile(policy *models.Policy, driverName string) ([]string, error) {
	driver, ok := c.drivers[driverName]
	if !ok {
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}

	var rules []string
	for _, rule := range policy.Rules {
		if !rule.Enabled {
			continue
		}
		compiled, err := driver.CompileRule(&rule)
		if err != nil {
			return nil, fmt.Errorf("failed to compile rule %q: %w", rule.Name, err)
		}
		rules = append(rules, compiled)
	}
	return rules, nil
}

func (c *Compiler) RegisterDriver(name string, driver drivers.Driver) {
	c.drivers[name] = driver
}
