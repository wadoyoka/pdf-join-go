package deleter

import (
	"fmt"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// DeletePages removes the specified pages from a PDF and writes the result.
// If outFile is empty, the input file is modified in place.
func DeletePages(inFile, outFile string, pages []int) error {
	totalPages, err := api.PageCountFile(inFile)
	if err != nil {
		return fmt.Errorf("failed to get page count: %w", err)
	}

	for _, p := range pages {
		if p < 1 || p > totalPages {
			return fmt.Errorf("page %d is out of range (1-%d)", p, totalPages)
		}
	}

	if len(pages) >= totalPages {
		return fmt.Errorf("cannot delete all %d pages", totalPages)
	}

	pageStrs := make([]string, len(pages))
	for i, p := range pages {
		pageStrs[i] = strconv.Itoa(p)
	}

	if err := api.RemovePagesFile(inFile, outFile, pageStrs, nil); err != nil {
		return fmt.Errorf("failed to delete pages: %w", err)
	}

	return nil
}
