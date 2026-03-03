package merger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectPDFs_SortedOrder(t *testing.T) {
	dir := t.TempDir()

	// Create files in non-alphabetical order
	for _, name := range []string{"c.pdf", "a.pdf", "b.pdf"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("dummy"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := CollectPDFs(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}

	expected := []string{
		filepath.Join(dir, "a.pdf"),
		filepath.Join(dir, "b.pdf"),
		filepath.Join(dir, "c.pdf"),
	}
	for i, f := range files {
		if f != expected[i] {
			t.Errorf("files[%d] = %s, want %s", i, f, expected[i])
		}
	}
}

func TestCollectPDFs_ExcludesNonPDF(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.pdf", "b.txt", "c.PDF", "d.png"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("dummy"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := CollectPDFs(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 PDF files, got %d: %v", len(files), files)
	}
}

func TestCollectPDFs_ExcludesDirectories(t *testing.T) {
	dir := t.TempDir()

	if err := os.Mkdir(filepath.Join(dir, "subdir.pdf"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.pdf"), []byte("dummy"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := CollectPDFs(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(files), files)
	}
}

func TestCollectPDFs_NonExistentDir(t *testing.T) {
	_, err := CollectPDFs("/nonexistent/dir", false)
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}

func TestCollectPDFs_Recursive(t *testing.T) {
	dir := t.TempDir()

	// Create nested structure:
	// dir/a.pdf
	// dir/sub1/b.pdf
	// dir/sub1/sub2/c.pdf
	// dir/other.txt  (should be excluded)
	if err := os.MkdirAll(filepath.Join(dir, "sub1", "sub2"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"a.pdf", "sub1/b.pdf", "sub1/sub2/c.pdf", "other.txt"} {
		if err := os.WriteFile(filepath.Join(dir, rel), []byte("dummy"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Non-recursive: only top-level
	flat, err := CollectPDFs(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(flat) != 1 {
		t.Fatalf("flat: expected 1 file, got %d: %v", len(flat), flat)
	}

	// Recursive: all PDFs
	rec, err := CollectPDFs(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec) != 3 {
		t.Fatalf("recursive: expected 3 files, got %d: %v", len(rec), rec)
	}

	// Verify sorted order
	expected := []string{
		filepath.Join(dir, "a.pdf"),
		filepath.Join(dir, "sub1", "b.pdf"),
		filepath.Join(dir, "sub1", "sub2", "c.pdf"),
	}
	for i, f := range rec {
		if f != expected[i] {
			t.Errorf("rec[%d] = %s, want %s", i, f, expected[i])
		}
	}
}

func TestCollectPDFs_RecursiveNonExistentDir(t *testing.T) {
	_, err := CollectPDFs("/nonexistent/dir", true)
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}

func TestMerge_TooFewFiles(t *testing.T) {
	err := Merge(nil, "out.pdf")
	if err == nil {
		t.Fatal("expected error for nil files")
	}

	err = Merge([]string{"a.pdf"}, "out.pdf")
	if err == nil {
		t.Fatal("expected error for single file")
	}
}
