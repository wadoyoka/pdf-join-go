package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:embed credits.txt
var credits string

var version = "dev"

var showCredits bool

var rootCmd = &cobra.Command{
	Use:     "nigopdf",
	Short:   "A PDF toolkit - merge and split PDF files",
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showCredits {
			fmt.Println(credits)
			return nil
		}
		return cmd.Help()
	},
}

func init() {
	rootCmd.Flags().BoolVar(&showCredits, "credits", false, "show third-party license information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
