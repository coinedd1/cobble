package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/coinedd1/cobble/internal/k8s"
	"github.com/coinedd1/cobble/internal/mc"
)

var sparkURLRe = regexp.MustCompile(`https?://spark\.lucko\.me/[A-Za-z0-9]+`)

var sparkWait time.Duration

var sparkCmd = &cobra.Command{
	Use:   "spark [args...]",
	Short: "Run a spark command and show the log output it produces",
	Args:  cobra.ArbitraryArgs,
	RunE: func(c *cobra.Command, args []string) error {
		conn, err := requireRCON()
		if err != nil {
			return err
		}

		line := "spark" // bare `cobble spark` → spark's help
		if len(args) > 0 {
			line = "spark " + strings.Join(args, " ")
		}

		// Mark the log position, backdated a little to absorb clock skew
		// between this machine and the node, then fire the command.
		since := time.Now().Add(-2 * time.Second)
		if _, err := conn.Command(line); err != nil {
			return err
		}

		// spark writes its output — including async upload URLs — to the
		// log, not the RCON reply. Wait briefly, then show what's new.
		time.Sleep(sparkWait)

		pod, err := requirePod()
		if err != nil {
			return err
		}
		out, err := k8s.LogsSince(cfg.Context, cfg.Namespace, pod, cfg.Container, since)
		if err != nil {
			return err
		}

		out = strings.TrimSpace(mc.Strip(out))
		if out == "" {
			fmt.Println("(no new log output — try a longer --wait if this was an upload)")
			return nil
		}
		fmt.Println(out)

		// Surface any result URL on its own line for easy copying.
		if url := sparkURLRe.FindString(out); url != "" {
			fmt.Println("\nspark:", url)
		}
		return nil
	},
}

func init() {
	sparkCmd.Flags().DurationVarP(&sparkWait, "wait", "w", 3*time.Second,
		"how long to wait for spark's log output before reading")
	rootCmd.AddCommand(sparkCmd)
}
