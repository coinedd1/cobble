package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/coinedd1/cobble/internal/k8s"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Copy files out of the server pod",
}

var pullCrashCmd = &cobra.Command{
	Use:   "crash [n]",
	Short: "Pull the most recent crash report (or the nth most recent)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		idx := 0
		if len(args) == 1 {
			n, err := strconv.Atoi(args[0])
			if err != nil || n < 1 {
				return fmt.Errorf("n must be a positive number (got %q)", args[0])
			}
			idx = n - 1
		}

		pod, err := requirePod()
		if err != nil {
			return err
		}

		dir := cfg.DataDir + "/crash-reports"
		files, err := k8s.ListFiles(cfg.Context, cfg.Namespace, pod, cfg.Container, dir)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return fmt.Errorf("no crash reports in %s (good news!)", dir)
		}
		if idx >= len(files) {
			return fmt.Errorf("only %d crash report(s) available", len(files))
		}

		chosen := files[idx]
		local := filepath.Base(chosen.Path)
		if err := k8s.Copy(cfg.Context, cfg.Namespace, pod, cfg.Container, chosen.Path, local); err != nil {
			return err
		}

		fmt.Printf("pulled %s (%s) → ./%s\n", filepath.Base(chosen.Path), humanAge(chosen.ModTime), local)
		return nil
	},
}

var pullLogCmd = &cobra.Command{
	Use:   "log [n]",
	Short: "Pull the most recent logs (or the nth most recent)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		idx := 0
		if len(args) == 1 {
			n, err := strconv.Atoi(args[0])
			if err != nil || n < 1 {
				return fmt.Errorf("n must be a positive number (got %q)", args[0])
			}
			idx = n - 1
		}

		pod, err := requirePod()
		if err != nil {
			return err
		}

		dir := cfg.DataDir + "/logs"
		files, err := k8s.ListFiles(cfg.Context, cfg.Namespace, pod, cfg.Container, dir)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return fmt.Errorf("no logs in %s", dir)
		}
		if idx >= len(files) {
			return fmt.Errorf("only %d log(s) available", len(files))
		}

		chosen := files[idx]
		local := filepath.Base(chosen.Path)
		if err := k8s.Copy(cfg.Context, cfg.Namespace, pod, cfg.Container, chosen.Path, local); err != nil {
			return err
		}

		fmt.Printf("pulled %s (%s) → ./%s\n", filepath.Base(chosen.Path), humanAge(chosen.ModTime), local)
		return nil
	},
}

var pullListCmd = &cobra.Command{
	Use:   "list [subdir]",
	Short: "List files available to pull (default: crash-reports), newest first",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		sub := "crash-reports"
		if len(args) == 1 {
			sub = args[0]
		}

		pod, err := requirePod()
		if err != nil {
			return err
		}

		dir := cfg.DataDir + "/" + sub
		files, err := k8s.ListFiles(cfg.Context, cfg.Namespace, pod, cfg.Container, dir)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Printf("no files in %s\n", dir)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		fmt.Fprintln(w, "  #\tAGE\tSIZE\tNAME")
		for i, f := range files {
			fmt.Fprintf(w, "  %d\t%s\t%s\t%s\n",
				i+1, humanAge(f.ModTime), humanSize(f.Size), filepath.Base(f.Path))
		}
		return w.Flush()
	},
}

func humanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGT"[exp])
}

func humanAge(t time.Time) string {
	d := time.Since(t).Round(time.Second)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func init() {
	pullCmd.AddCommand(pullCrashCmd, pullLogCmd, pullListCmd)
	rootCmd.AddCommand(pullCmd)
}
