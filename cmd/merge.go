package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wadoyoka/nigopdf/internal/merger"
)

var mergeOutput string
var mergeDryRun bool
var mergeRecursive bool

var mergeCmd = &cobra.Command{
	Use:   "merge [directory]",
	Short: "Merge PDF files in a directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		files, err := merger.CollectPDFs(dir, mergeRecursive)
		if err != nil {
			return err
		}

		if mergeDryRun {
			if len(files) == 0 {
				fmt.Println("No PDF files found")
				return nil
			}
			for _, f := range files {
				fmt.Println(f)
			}
			return nil
		}

		if err := merger.Merge(files, mergeOutput); err != nil {
			return err
		}

		fmt.Printf("Merged %d files into %s\n", len(files), mergeOutput)
		return nil
	},
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "merged.pdf", "output file name")
	mergeCmd.Flags().BoolVarP(&mergeDryRun, "dry-run", "n", false, "list target files without merging")
	mergeCmd.Flags().BoolVarP(&mergeRecursive, "recursive", "r", false, "search subdirectories recursively")
	rootCmd.AddCommand(mergeCmd)
}
