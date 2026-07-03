package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/coinedd1/cobble/internal/config"
	"github.com/coinedd1/cobble/internal/k8s"
	"github.com/coinedd1/cobble/internal/mc"
	"github.com/coinedd1/cobble/internal/rcon"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:           "cobble",
	Short:         "Manage a Minecraft server running on Kubernetes",
	Long:          "cobble consolidates normally tedious k8s management of a Minecraft server into a CLI tool that supports file pulls, in-game command execution, world management, and more.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		return nil
	},
}

var resolvedPod string

func requirePod() (string, error) {
	if resolvedPod != "" {
		return resolvedPod, nil
	}
	pod, err := k8s.ResolvePod(cfg.Context, cfg.Namespace, cfg.Selector)
	if err != nil {
		return "", err
	}
	resolvedPod = pod
	return pod, nil
}

func requireRCON() (rcon.RCON, error) {
	pod, err := requirePod()
	if err != nil {
		return nil, err
	}
	return rcon.ExecRCON{
		Context:   cfg.Context,
		Namespace: cfg.Namespace,
		Pod:       pod,
		Container: cfg.Container,
	}, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func send(line string) error {
	conn, err := requireRCON()
	if err != nil {
		return err
	}
	out, err := conn.Command(line)
	if err != nil {
		return err
	}
	if out != "" {
		fmt.Println(mc.ToANSI(out))
	}
	return nil
}
