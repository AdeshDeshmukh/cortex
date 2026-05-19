package commands

import (
	"github.com/spf13/cobra"
)

var Version = "dev"

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "cortex",
	Short: "Cortex - Your code's thinking layer",
	Long: `Cortex is a privacy-first AI code reviewer that runs locally
and adapts to your coding style through reinforcement learning.

Built from first principles to understand code deeply.

Example usage:
  cortex install          Install git hooks in current repo
  cortex review           Review current changes`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .cortex.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}