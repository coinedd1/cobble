package k8s

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ResolvePod(ctx, namespace, selector string) (string, error) {
	args := []string{}
	if ctx != "" {
		args = append(args, "--context", ctx)
	}
	args = append(
		args,
		"get", "pod",
		"-n", namespace,
		"-l", selector,
		"-o", "jsonpath={range .items[*]}{.metadata.name}{'\\n'}{end}",
	)

	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) && len(ee.Stderr) > 0 {
			return "", fmt.Errorf("resolving pod: %s", strings.TrimSpace(string(ee.Stderr)))
		}
		return "", fmt.Errorf("resolving pod: %w", err)
	}
	name := strings.TrimSpace(string(out))
	if name == "" {
		return "", fmt.Errorf("no pod matches selector %q in namespace %q", selector, namespace)
	}
	return name, nil
}

func Logs(ctx, namespace, pod, container string, tail int) (string, error) {
	args := []string{}
	if ctx != "" {
		args = append(args, "--context", ctx)
	}
	args = append(args, "logs", "-n", namespace, pod)
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, fmt.Sprintf("--tail=%d", tail))

	out, err := exec.Command("kubectl", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("logs: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func LogsSince(ctx, namespace, pod, container string, since time.Time) (string, error) {
	args := []string{}
	if ctx != "" {
		args = append(args, "--context", ctx)
	}
	args = append(args, "logs", "-n", namespace, pod)
	if container != "" {
		args = append(args, "-c", container)
	}
	// --since-time wants RFC3339 (second precision).
	args = append(args, "--since-time="+since.UTC().Format(time.RFC3339))

	out, err := exec.Command("kubectl", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("logs: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
