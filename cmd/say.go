package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/coinedd1/cobble/internal/mc"
)

var sayCmd = &cobra.Command{
	Use:   "say <message>",
	Short: "Broadcast a message and preview how it renders",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		msg := mc.Amp(strings.Join(args, " ")) // &a → §a before sending

		conn, err := requireRCON()
		if err != nil {
			return err
		}
		if _, err := conn.Command("say " + msg); err != nil {
			return err
		}

		// Reply is empty, so we render the sent message as the "output".
		fmt.Println("sent:", mc.ToANSI(msg))
		return nil
	},
}

func init() { rootCmd.AddCommand(sayCmd) }
