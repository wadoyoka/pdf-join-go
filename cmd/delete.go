package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wadoyoka/nigopdf/internal/deleter"
	"github.com/wadoyoka/nigopdf/internal/pageutil"
)

var deletePages string
var deleteOutput string

var deleteCmd = &cobra.Command{
	Use:   "delete <file>",
	Short: "Delete specific pages from a PDF file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pages, err := pageutil.ParsePages(deletePages)
		if err != nil {
			return err
		}

		inFile := args[0]

		outFile := deleteOutput
		if outFile == "" {
			ext := filepath.Ext(inFile)
			base := strings.TrimSuffix(inFile, ext)
			outFile = base + "_deleted" + ext
		}

		if err := deleter.DeletePages(inFile, outFile, pages); err != nil {
			return err
		}

		fmt.Printf("Deleted pages %s, saved to %s\n", deletePages, outFile)
		return nil
	},
}

func init() {
	deleteCmd.Flags().StringVar(&deletePages, "pages", "", "comma-separated page numbers to delete (e.g. 2,5,8)")
	_ = deleteCmd.MarkFlagRequired("pages")
	deleteCmd.Flags().StringVarP(&deleteOutput, "output", "o", "", "output file (default: <input>_deleted.pdf)")
	rootCmd.AddCommand(deleteCmd)
}
