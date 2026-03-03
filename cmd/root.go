package cmd

import (
	"fmt"
	"os"

	"github.com/wadoyoka/pdf-join-go/internal/merger"
	"github.com/spf13/cobra"
)

var output string

var rootCmd = &cobra.Command{
	Use:   "pdf-join [directory]",
	Short: "Merge PDF files in a directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
