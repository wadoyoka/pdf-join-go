package deleter

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

// createTestPDF creates a PDF with the given number of pages at path.
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

func TestDeletePages_Basic(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "input.pdf")
	createTestPDF(t, inFile, 10)

	err := DeletePages(inFile, "", []int{2, 5, 8})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count, err := api.PageCountFile(inFile)
	if err != nil {
		t.Fatalf("failed to get page count: %v", err)
	}
	if count != 7 {
		t.Errorf("expected 7 pages, got %d", count)
	}
}

func TestDeletePages_FirstPage(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "input.pdf")
	createTestPDF(t, inFile, 5)

	err := DeletePages(inFile, "", []int{1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count, err := api.PageCountFile(inFile)
	if err != nil {
		t.Fatalf("failed to get page count: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 pages, got %d", count)
	}
}

func TestDeletePages_LastPage(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "input.pdf")
	createTestPDF(t, inFile, 5)

	err := DeletePages(inFile, "", []int{5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count, err := api.PageCountFile(inFile)
	if err != nil {
		t.Fatalf("failed to get page count: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 pages, got %d", count)
	}
}

func TestDeletePages_OutOfRange(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "input.pdf")
	createTestPDF(t, inFile, 5)

	err := DeletePages(inFile, "", []int{6})
	if err == nil {
		t.Fatal("expected error for out-of-range page, got nil")
	}
}

func TestDeletePages_AllPages(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "input.pdf")
	createTestPDF(t, inFile, 3)

	err := DeletePages(inFile, "", []int{1, 2, 3})
	if err == nil {
		t.Fatal("expected error when deleting all pages, got nil")
	}
}

func TestDeletePages_WithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inFile := filepath.Join(tmpDir, "input.pdf")
	outFile := filepath.Join(tmpDir, "output.pdf")
	createTestPDF(t, inFile, 5)

	err := DeletePages(inFile, outFile, []int{2, 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Input should remain unchanged at 5 pages.
	inCount, err := api.PageCountFile(inFile)
	if err != nil {
		t.Fatalf("failed to get input page count: %v", err)
	}
	if inCount != 5 {
		t.Errorf("expected input to have 5 pages, got %d", inCount)
	}

	// Output should have 3 pages.
	outCount, err := api.PageCountFile(outFile)
	if err != nil {
		t.Fatalf("failed to get output page count: %v", err)
	}
	if outCount != 3 {
		t.Errorf("expected output to have 3 pages, got %d", outCount)
	}

	// Output file should exist.
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}
