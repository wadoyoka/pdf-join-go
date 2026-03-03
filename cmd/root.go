package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wadoyoka/pdf-join-go/internal/merger"
)

//go:embed credits.txt
var credits string

var output string
var showCredits bool

var rootCmd = &cobra.Command{
	Use:   "pdf-join [directory]",
	Short: "Merge PDF files in a directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if showCredits {
			fmt.Println(credits)
			return nil
		}

		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		files, err := merger.CollectPDFs(dir)
		if err != nil {
			return err
		}

		if err := merger.Merge(files, output); err != nil {
			return err
		}

		fmt.Printf("Merged %d files into %s\n", len(files), output)
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&output, "output", "o", "merged.pdf", "output file name")
	rootCmd.Flags().BoolVar(&showCredits, "credits", false, "show third-party license information")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
