package splitter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// createSinglePagePDF creates a valid single-page PDF file at path.
func createSinglePagePDF(t *testing.T, path string) {
	t.Helper()
	xRefTable, err := pdfcpu.CreateXRefTableWithRootDict()
	if err != nil {
		t.Fatalf("failed to create xref table: %v", err)
	}
	rootDict, err := xRefTable.Catalog()
	if err != nil {
		t.Fatalf("failed to get catalog: %v", err)
	}
	mediaBox := types.RectForFormat("A4")
	p := model.Page{MediaBox: mediaBox, Buf: new(bytes.Buffer)}
	if err := pdfcpu.AddPageTreeWithSamplePage(xRefTable, rootDict, p); err != nil {
		t.Fatalf("failed to add page tree: %v", err)
	}
	ctx := pdfcpu.CreateContext(xRefTable, model.NewDefaultConfiguration())
	if err := api.WriteContextFile(ctx, path); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}
}

// createTestPDF generates a valid PDF with the given number of pages.
func createTestPDF(t *testing.T, path string, pages int) {
	t.Helper()
	if pages < 1 {
		t.Fatal("pages must be >= 1")
	}

	tmpDir := t.TempDir()

	if pages == 1 {
		createSinglePagePDF(t, path)
		return
	}

	var files []string
	for i := range pages {
		f := filepath.Join(tmpDir, fmt.Sprintf("p%d.pdf", i))
		createSinglePagePDF(t, f)
		files = append(files, f)
	}

	if err := api.MergeCreateFile(files, path, false, nil); err != nil {
		t.Fatalf("failed to merge test PDF: %v", err)
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantErr bool
	}{
		{"10MB", 10 * 1024 * 1024, false},
		{"500KB", 500 * 1024, false},
		{"1GB", 1024 * 1024 * 1024, false},
		{"100B", 100, false},
		{"10mb", 10 * 1024 * 1024, false},
		{"1.5MB", int64(1.5 * 1024 * 1024), false},
		{"abc", 0, true},
		{"", 0, true},
		{"-10MB", 0, true},
		{"0MB", 0, true},
		{"MB", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseSize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseSize(%q) = %d, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseSize(%q) error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSize(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestOutputFileName(t *testing.T) {
	tests := []struct {
		outDir   string
		basename string
		index    int
		want     string
	}{
		{"/out", "doc", 1, "/out/doc_1.pdf"},
		{"/tmp", "report", 3, "/tmp/report_3.pdf"},
		{".", "test", 10, "test_10.pdf"},
	}

	for _, tt := range tests {
		got := outputFileName(tt.outDir, tt.basename, tt.index)
		if got != tt.want {
			t.Errorf("outputFileName(%q, %q, %d) = %q, want %q", tt.outDir, tt.basename, tt.index, got, tt.want)
		}
	}
}

func TestSplitByParts(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "input.pdf")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatal(err)
	}

	createTestPDF(t, inFile, 6)

	files, err := SplitByParts(inFile, outDir, 3)
	if err != nil {
		t.Fatalf("SplitByParts failed: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}

	// Verify each output file exists and is a valid PDF
	for i, f := range files {
		expected := filepath.Join(outDir, fmt.Sprintf("input_%d.pdf", i+1))
		if f != expected {
			t.Errorf("file[%d] = %q, want %q", i, f, expected)
		}

		count, err := api.PageCountFile(f)
		if err != nil {
			t.Errorf("failed to read page count for %s: %v", f, err)
			continue
		}
		if count != 2 {
			t.Errorf("part %d has %d pages, want 2", i+1, count)
		}
	}
}

func TestSplitByParts_UnevenPages(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "input.pdf")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatal(err)
	}

	createTestPDF(t, inFile, 7)

	files, err := SplitByParts(inFile, outDir, 3)
	if err != nil {
		t.Fatalf("SplitByParts failed: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}

	// 7 pages / 3 parts = 2 base + 1 remainder
	// First part gets 3 pages, rest get 2 each
	expectedPages := []int{3, 2, 2}
	for i, f := range files {
		count, err := api.PageCountFile(f)
		if err != nil {
			t.Errorf("failed to read page count for %s: %v", f, err)
			continue
		}
		if count != expectedPages[i] {
			t.Errorf("part %d has %d pages, want %d", i+1, count, expectedPages[i])
		}
	}
}

func TestSplitByParts_ExceedingPages(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "input.pdf")
	createTestPDF(t, inFile, 3)

	_, err := SplitByParts(inFile, dir, 5)
	if err == nil {
		t.Fatal("expected error when parts > pages")
	}
}

func TestSplitByMaxSize(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "input.pdf")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatal(err)
	}

	createTestPDF(t, inFile, 4)

	// Use a very large max size to get a single output file
	files, err := SplitByMaxSize(inFile, outDir, 100*1024*1024)
	if err != nil {
		t.Fatalf("SplitByMaxSize failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file with large max-size, got %d", len(files))
	}

	count, err := api.PageCountFile(files[0])
	if err != nil {
		t.Fatalf("failed to read page count: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 pages, got %d", count)
	}
}

func TestSplitByMaxSize_SmallSize(t *testing.T) {
	dir := t.TempDir()
	inFile := filepath.Join(dir, "input.pdf")
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatal(err)
	}

	createTestPDF(t, inFile, 4)

	// Get a single page's size to use as max-size
	info, err := os.Stat(inFile)
	if err != nil {
		t.Fatal(err)
	}
	// Use half the file size as max-size to force multiple parts
	maxSize := info.Size() / 2
	if maxSize < 1024 {
		maxSize = 1024
	}

	files, err := SplitByMaxSize(inFile, outDir, maxSize)
	if err != nil {
		t.Fatalf("SplitByMaxSize failed: %v", err)
	}

	if len(files) < 2 {
		t.Fatalf("expected at least 2 files with small max-size, got %d", len(files))
	}

	// Verify total pages across all parts equals original
	totalPages := 0
	for _, f := range files {
		count, err := api.PageCountFile(f)
		if err != nil {
			t.Errorf("failed to read page count for %s: %v", f, err)
			continue
		}
		totalPages += count
	}
	if totalPages != 4 {
		t.Errorf("total pages across parts = %d, want 4", totalPages)
	}
}
