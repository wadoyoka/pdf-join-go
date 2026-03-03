package splitter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

var sizePattern = regexp.MustCompile(`(?i)^(\d+(?:\.\d+)?)\s*(B|KB|MB|GB)$`)

// ParseSize parses a human-readable size string (e.g. "10MB", "500KB") into bytes.
func ParseSize(s string) (int64, error) {
	m := sizePattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, fmt.Errorf("invalid size format: %q (expected e.g. 10MB, 500KB)", s)
	}

	val, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size number: %w", err)
	}
	if val <= 0 {
		return 0, fmt.Errorf("size must be positive, got %s", s)
	}

	var multiplier float64
	switch strings.ToUpper(m[2]) {
	case "B":
		multiplier = 1
	case "KB":
		multiplier = 1024
	case "MB":
		multiplier = 1024 * 1024
	case "GB":
		multiplier = 1024 * 1024 * 1024
	}

	return int64(val * multiplier), nil
}

// outputFileName returns the output path for a split part.
// e.g. outDir="/out", basename="doc", index=1 -> "/out/doc_1.pdf"
func outputFileName(outDir, basename string, index int) string {
	return filepath.Join(outDir, fmt.Sprintf("%s_%d.pdf", basename, index))
}

// SplitByParts splits a PDF into the given number of parts with roughly equal page counts.
// Extra pages (remainder) are distributed one per part to the first parts.
func SplitByParts(inFile, outDir string, parts int) ([]string, error) {
	totalPages, err := api.PageCountFile(inFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}
	if parts > totalPages {
		return nil, fmt.Errorf("cannot split %d pages into %d parts", totalPages, parts)
	}

	base := totalPages / parts
	remainder := totalPages % parts

	f, err := os.Open(inFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	basename := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))
	var outFiles []string

	pageStart := 1
	for i := 0; i < parts; i++ {
		count := base
		if i < remainder {
			count++
		}
		pageEnd := pageStart + count - 1

		pages := api.PagesForPageRange(pageStart, pageEnd)

		ctx, err := readContext(f)
		if err != nil {
			return nil, err
		}

		extracted, err := pdfcpu.ExtractPages(ctx, pages, false)
		if err != nil {
			return nil, fmt.Errorf("failed to extract pages %d-%d: %w", pageStart, pageEnd, err)
		}

		outFile := outputFileName(outDir, basename, i+1)
		if err := api.WriteContextFile(extracted, outFile); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", outFile, err)
		}

		outFiles = append(outFiles, outFile)
		pageStart = pageEnd + 1
	}

	return outFiles, nil
}

// SplitByMaxSize splits a PDF so that each output part stays under maxBytes.
// Individual page sizes are estimated by splitting into single pages. Due to shared
// resources (fonts, etc.), actual part sizes will typically be smaller than maxBytes.
func SplitByMaxSize(inFile, outDir string, maxBytes int64) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	spans, err := api.SplitRaw(f, 1, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to split into pages: %w", err)
	}

	// Measure each page's size
	pageSizes := make([]int64, len(spans))
	for i, span := range spans {
		data, err := io.ReadAll(span.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read page %d: %w", i+1, err)
		}
		pageSizes[i] = int64(len(data))

		if pageSizes[i] > maxBytes {
			return nil, fmt.Errorf("page %d (%d bytes) exceeds max-size (%d bytes)", i+1, pageSizes[i], maxBytes)
		}
	}

	// Greedy grouping
	var groups [][]int // each group is a list of 1-based page numbers
	var currentGroup []int
	var currentSize int64

	for i, size := range pageSizes {
		pageNum := i + 1
		if len(currentGroup) > 0 && currentSize+size > maxBytes {
			groups = append(groups, currentGroup)
			currentGroup = nil
			currentSize = 0
		}
		currentGroup = append(currentGroup, pageNum)
		currentSize += size
	}
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	// Write each group
	basename := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))
	var outFiles []string

	for i, group := range groups {
		ctx, err := readContext(f)
		if err != nil {
			return nil, err
		}

		extracted, err := pdfcpu.ExtractPages(ctx, group, false)
		if err != nil {
			return nil, fmt.Errorf("failed to extract pages for part %d: %w", i+1, err)
		}

		outFile := outputFileName(outDir, basename, i+1)
		if err := api.WriteContextFile(extracted, outFile); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", outFile, err)
		}

		outFiles = append(outFiles, outFile)
	}

	return outFiles, nil
}

func readContext(f *os.File) (*model.Context, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}
	conf := model.NewDefaultConfiguration()
	ctx, err := api.ReadValidateAndOptimize(f, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}
	return ctx, nil
}
