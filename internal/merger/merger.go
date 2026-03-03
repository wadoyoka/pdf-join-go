package merger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// CollectPDFs returns .pdf files in the directory sorted by filename ascending.
func CollectPDFs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".pdf") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

// Merge combines the given PDF files into outFile.
func Merge(files []string, outFile string) error {
	if len(files) < 2 {
		return fmt.Errorf("need at least 2 PDF files to merge, found %d", len(files))
	}
	return api.MergeCreateFile(files, outFile, false, nil)
}
