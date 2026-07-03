package rcon

import (
	"fmt"
	"os/exec"
	"strings"
)

type RCON interface {
	Command(line string) (string, error)
}

type ExecRCON struct {
	Context   string
	Namespace string
	Pod       string
	Container string
}

func (e ExecRCON) Command(line string) (string, error) {
	args := []string{}
	if e.Context != "" {
		args = append(args, "--context", e.Context)
	}
	args = append(
		args,
		"exec", "-n", e.Namespace, e.Pod,
	)
	if e.Container != "" {
		args = append(args, "-c", e.Container)
	}
	args = append(args, "--", "rcon-cli", line)

	out, err := exec.Command("kubectl", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("rcon-cli: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}
