package k8s

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
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

func Exec(ctx, namespace, pod, container string, command ...string) (string, error) {
	args := []string{}
	if ctx != "" {
		args = append(args, "--context", ctx)
	}
	args = append(args, "exec", "-n", namespace, pod)
	if container != "" {
		args = append(args, "-c", container)
	}
	args = append(args, "--")
	args = append(args, command...)

	out, err := exec.Command("kubectl", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

type RemoteFile struct {
	ModTime time.Time
	Size    int64
	Path    string
}

func ListFiles(ctx, namespace, pod, container, dir string) ([]RemoteFile, error) {
	script := fmt.Sprintf(
		"find '%s' -type f -printf '%%T@\\t%%s\\t%%p\\n' 2>/dev/null | sort -rn",
		dir,
	)
	out, err := Exec(ctx, namespace, pod, container, "sh", "-c", script)
	if err != nil {
		return nil, err
	}

	var files []RemoteFile
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			continue
		}
		secs, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			continue
		}
		size, _ := strconv.ParseInt(parts[1], 10, 64)
		whole := int64(secs)
		frac := int64((secs - float64(whole)) * 1e9)
		files = append(files, RemoteFile{
			ModTime: time.Unix(whole, frac),
			Size:    size,
			Path:    parts[2],
		})
	}
	return files, nil
}

func Copy(ctx, namespace, pod, container, remotePath, localPath string) error {
	src := fmt.Sprintf("%s/%s:%s", namespace, pod, remotePath)
	args := []string{}
	if ctx != "" {
		args = append(args, "--context", ctx)
	}
	args = append(args, "cp", src, localPath)
	if container != "" {
		args = append(args, "-c", container)
	}

	out, err := exec.Command("kubectl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("cp: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
