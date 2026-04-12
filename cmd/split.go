package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wadoyoka/nigopdf/internal/pageutil"
	"github.com/wadoyoka/nigopdf/internal/splitter"
)

var splitParts int
var splitMaxSize string
var splitPages string
var splitOutput string

var splitCmd = &cobra.Command{
	Use:   "split <file>",
	Short: "Split a PDF file into multiple parts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		partsSet := cmd.Flags().Changed("parts")
		maxSizeSet := cmd.Flags().Changed("max-size")
		pagesSet := cmd.Flags().Changed("pages")

		setCount := 0
		if partsSet {
			setCount++
		}
		if maxSizeSet {
			setCount++
		}
		if pagesSet {
			setCount++
		}

		if setCount > 1 {
			return fmt.Errorf("--parts, --max-size, and --pages are mutually exclusive")
		}
		if setCount == 0 {
			return fmt.Errorf("specify one of --parts, --max-size, or --pages")
		}

		inFile := args[0]

		outDir := splitOutput
		if outDir == "" {
			outDir = filepath.Dir(inFile)
		}

		if pagesSet {
			pages, err := pageutil.ParsePages(splitPages)
			if err != nil {
				return err
			}
			files, err := splitter.SplitByPages(inFile, outDir, pages)
			if err != nil {
				return err
			}
			for _, f := range files {
				fmt.Println(f)
			}
			fmt.Printf("Split into %d parts\n", len(files))
			return nil
		}

		if partsSet {
			if splitParts < 2 {
				return fmt.Errorf("--parts must be at least 2, got %d", splitParts)
			}
			files, err := splitter.SplitByParts(inFile, outDir, splitParts)
			if err != nil {
				return err
			}
			for _, f := range files {
				fmt.Println(f)
			}
			fmt.Printf("Split into %d parts\n", len(files))
			return nil
		}

		maxBytes, err := splitter.ParseSize(splitMaxSize)
		if err != nil {
			return err
		}
		files, err := splitter.SplitByMaxSize(inFile, outDir, maxBytes)
		if err != nil {
			return err
		}
		for _, f := range files {
			fmt.Println(f)
		}
		fmt.Printf("Split into %d parts\n", len(files))
		return nil
	},
}

func init() {
	splitCmd.Flags().IntVar(&splitParts, "parts", 0, "number of parts to split into")
	splitCmd.Flags().StringVar(&splitMaxSize, "max-size", "", "maximum size per part (e.g. 10MB, 500KB)")
	splitCmd.Flags().StringVar(&splitPages, "pages", "", "page boundaries to split at (e.g. 2,5,8)")
	splitCmd.Flags().StringVarP(&splitOutput, "output", "o", "", "output directory (default: same as input file)")
	rootCmd.AddCommand(splitCmd)
}
