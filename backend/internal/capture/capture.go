package capture

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type FirewallState struct {
	Tables    []Table   `json:"tables"`
	Timestamp time.Time `json:"timestamp"`
}

type Table struct {
	Name   string           `json:"name"`
	Chains map[string]Chain `json:"chains"`
}

type Chain struct {
	Policy  string `json:"policy"`
	Rules   []Rule `json:"rules"`
	Packets uint64 `json:"packets,omitempty"`
	Bytes   uint64 `json:"bytes,omitempty"`
}

type Rule struct {
	Chain       string   `json:"chain"`
	Target      string   `json:"target"`
	Protocol    string   `json:"protocol,omitempty"`
	Source      string   `json:"source,omitempty"`
	Destination string   `json:"destination,omitempty"`
	SrcPort     string   `json:"src_port,omitempty"`
	DstPort     string   `json:"dst_port,omitempty"`
	InInterface string   `json:"in_interface,omitempty"`
	OutInterface string  `json:"out_interface,omitempty"`
	State       string   `json:"state,omitempty"`
	LogPrefix   string   `json:"log_prefix,omitempty"`
	Extra       []string `json:"extra,omitempty"`
	Raw         string   `json:"raw"`
}

func Capture() (*FirewallState, error) {
	data, err := exec.Command("iptables-save", "-c").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run iptables-save: %w", err)
	}

	state, err := parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse iptables-save output: %w", err)
	}

	state.Timestamp = time.Now()
	return state, nil
}

func parse(data string) (*FirewallState, error) {
	state := &FirewallState{
		Timestamp: time.Now(),
	}

	var currentTable *Table

	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "*"):
			name := strings.TrimSpace(line[1:])
			currentTable = &Table{
				Name:   name,
				Chains: make(map[string]Chain),
			}
			state.Tables = append(state.Tables, *currentTable)

		case strings.HasPrefix(line, ":"):
			if currentTable == nil {
				continue
			}
			parts := strings.Fields(line[1:])
			if len(parts) < 2 {
				continue
			}
			ch := Chain{Policy: parts[1]}
			if len(parts) >= 3 {
				fmt.Sscanf(strings.Trim(parts[2], "[]"), "%d:%d", &ch.Packets, &ch.Bytes)
			}
			idx := len(state.Tables) - 1
			state.Tables[idx].Chains[parts[0]] = ch

		case strings.HasPrefix(line, "-A"):
			if currentTable == nil {
				continue
			}
			rule := parseRule(line)
			idx := len(state.Tables) - 1
			ch := state.Tables[idx].Chains[rule.Chain]
			ch.Rules = append(ch.Rules, rule)
			state.Tables[idx].Chains[rule.Chain] = ch

		case line == "COMMIT":
			currentTable = nil
		}
	}

	return state, nil
}

func splitTokens(line string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false

	for i := 0; i < len(line); i++ {
		c := line[i]
		if c == '"' {
			inQuote = !inQuote
			cur.WriteByte(c)
			continue
		}
		if c == ' ' && !inQuote {
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteByte(c)
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

func parseRule(line string) Rule {
	r := Rule{Raw: line}
	tokens := splitTokens(line)
	if len(tokens) < 2 {
		return r
	}
	r.Chain = tokens[1]

	for i := 2; i < len(tokens); i++ {
		t := tokens[i]
		if t == "!" {
			continue
		}
		switch t {
		case "-j":
			if i+1 < len(tokens) {
				r.Target = tokens[i+1]
				i++
			}
		case "-p":
			if i+1 < len(tokens) {
				r.Protocol = tokens[i+1]
				i++
			}
		case "-s":
			if i+1 < len(tokens) {
				r.Source = tokens[i+1]
				i++
			}
		case "-d":
			if i+1 < len(tokens) {
				r.Destination = tokens[i+1]
				i++
			}
		case "--sport":
			if i+1 < len(tokens) {
				r.SrcPort = tokens[i+1]
				i++
			}
		case "--dport":
			if i+1 < len(tokens) {
				r.DstPort = tokens[i+1]
				i++
			}
		case "-i":
			if i+1 < len(tokens) {
				r.InInterface = tokens[i+1]
				i++
			}
		case "-o":
			if i+1 < len(tokens) {
				r.OutInterface = tokens[i+1]
				i++
			}
		case "-m":
			if i+1 < len(tokens) {
				module := tokens[i+1]
				i++
				if module == "state" {
					for j := i + 1; j < len(tokens); j++ {
						if tokens[j] == "--state" && j+1 < len(tokens) {
							r.State = tokens[j+1]
							i = j + 1
							break
						}
						if strings.HasPrefix(tokens[j], "-") {
							i = j - 1
							break
						}
					}
				}
			}
		case "--log-prefix":
			if i+1 < len(tokens) {
				r.LogPrefix = strings.Trim(tokens[i+1], "\"")
				i++
			}
		default:
			if strings.HasPrefix(t, "-") {
				continue
			}
			r.Extra = append(r.Extra, t)
		}
	}

	return r
}
