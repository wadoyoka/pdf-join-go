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

	files, err := CollectPDFs(dir)
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

	files, err := CollectPDFs(dir)
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

	files, err := CollectPDFs(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(files), files)
	}
}

func TestCollectPDFs_NonExistentDir(t *testing.T) {
	_, err := CollectPDFs("/nonexistent/dir")
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
