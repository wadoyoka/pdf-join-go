package merger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func isPDF(name string) bool {
	return strings.EqualFold(filepath.Ext(name), ".pdf")
}

// CollectPDFs returns .pdf files in the directory sorted by path ascending.
// If recursive is true, subdirectories are searched as well.
func CollectPDFs(dir string, recursive bool) ([]string, error) {
	if recursive {
		return collectPDFsRecursive(dir)
	}
	return collectPDFsFlat(dir)
}

func collectPDFsFlat(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isPDF(e.Name()) {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

func collectPDFsRecursive(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if isPDF(d.Name()) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dir, err)
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
