package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd <console command>",
	Short: "Run a raw RCON console command",
	Long:  "Send a command to the server console over RCON and print the reply. \n\n Examples: cobble cmd list\n cobble cmd \"say hello everyone\"",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		conn, err := requireRCON()
		if err != nil {
			return err
		}

		line := strings.Join(args, " ")
		out, err := conn.Command(line)
		if err != nil {
			return err
		}
		if out != "" {
			fmt.Println(out)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cmdCmd)
}

var carpetPlayerCmd = &cobra.Command{
	Use:   "player <name> <spawn|kill|move|...>",
	Short: "Carpet mod player command",
	Long:  "Use carpet mod's player command. This creates a fake player in the server and allows you to perform basic functions with them",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(c *cobra.Command, args []string) error {
		return send("player " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(carpetPlayerCmd)
}

var kickCmd = &cobra.Command{
	Use:   "kick <player> <reason>",
	Short: "Kick a player",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return send("kick " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(kickCmd)
}

var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Manage the whitelist",
	// no RunE: running `cobble whitelist` alone prints the subcommand list
}

var whitelistAddCmd = &cobra.Command{
	Use:   "add <player>",
	Short: "Add a player to the whitelist",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return send("whitelist add " + args[0])
	},
}

var whitelistRemoveCmd = &cobra.Command{
	Use:   "remove <player>",
	Short: "Remove a player from the whitelist",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return send("whitelist remove " + args[0])
	},
}

var whitelistListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show whitelisted players",
	Args:  cobra.NoArgs,
	RunE: func(c *cobra.Command, args []string) error {
		return send("whitelist list")
	},
}

func init() {
	whitelistCmd.AddCommand(whitelistAddCmd, whitelistRemoveCmd, whitelistListCmd)
	rootCmd.AddCommand(whitelistCmd)
}

var opCmd = &cobra.Command{
	Use:   "op <player>",
	Short: "Give a player operator permissions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return send("op " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(opCmd)
}

var deopCmd = &cobra.Command{
	Use:   "deop <player>",
	Short: "Remove a player's operator permissions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return send("deop " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(deopCmd)
}
