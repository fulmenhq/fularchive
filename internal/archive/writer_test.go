package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fulmenhq/refbolt/internal/provider"
)

func TestWrite_NewFiles(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	pages := []provider.Page{
		{Path: "overview.md", Content: []byte("# Overview\n")},
		{Path: "guide/setup.md", Content: []byte("# Setup\n")},
	}

	stat, err := w.Write("llm-api", "test-provider", pages)
	if err != nil {
		t.Fatal(err)
	}
	if stat.Written != 2 {
		t.Errorf("Written = %d, want 2", stat.Written)
	}
	if stat.Skipped != 0 {
		t.Errorf("Skipped = %d, want 0", stat.Skipped)
	}
	if stat.Total != 2 {
		t.Errorf("Total = %d, want 2", stat.Total)
	}
}

func TestWrite_UnchangedFiles(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	pages := []provider.Page{
		{Path: "overview.md", Content: []byte("# Overview\n")},
	}

	// First write.
	stat1, err := w.Write("llm-api", "test-provider", pages)
	if err != nil {
		t.Fatal(err)
	}
	if stat1.Written != 1 {
		t.Fatalf("first write: Written = %d, want 1", stat1.Written)
	}

	// Second write with same content — should skip.
	stat2, err := w.Write("llm-api", "test-provider", pages)
	if err != nil {
		t.Fatal(err)
	}
	if stat2.Skipped != 1 {
		t.Errorf("second write: Skipped = %d, want 1", stat2.Skipped)
	}
	if stat2.Written != 0 {
		t.Errorf("second write: Written = %d, want 0", stat2.Written)
	}
}

func TestWrite_ChangedFile(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	// First write.
	pages1 := []provider.Page{
		{Path: "overview.md", Content: []byte("# Overview v1\n")},
	}
	if _, err := w.Write("llm-api", "test-provider", pages1); err != nil {
		t.Fatal(err)
	}

	// Second write with different content — should write.
	pages2 := []provider.Page{
		{Path: "overview.md", Content: []byte("# Overview v2 — updated\n")},
	}
	stat, err := w.Write("llm-api", "test-provider", pages2)
	if err != nil {
		t.Fatal(err)
	}
	if stat.Written != 1 {
		t.Errorf("Written = %d, want 1 (content changed)", stat.Written)
	}
	if stat.Skipped != 0 {
		t.Errorf("Skipped = %d, want 0", stat.Skipped)
	}

	// Verify new content was written.
	dateDir := findDateDir(t, filepath.Join(dir, "llm-api", "test-provider"))
	content, err := os.ReadFile(filepath.Join(dateDir, "overview.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# Overview v2 — updated\n" {
		t.Errorf("unexpected content: %q", content)
	}
}

func TestWrite_MixedBatch(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	// First write: 3 files.
	pages1 := []provider.Page{
		{Path: "a.md", Content: []byte("a")},
		{Path: "b.md", Content: []byte("b")},
		{Path: "c.md", Content: []byte("c")},
	}
	if _, err := w.Write("topic", "prov", pages1); err != nil {
		t.Fatal(err)
	}

	// Second write: a unchanged, b changed, d new.
	pages2 := []provider.Page{
		{Path: "a.md", Content: []byte("a")},         // unchanged
		{Path: "b.md", Content: []byte("b-updated")}, // changed
		{Path: "d.md", Content: []byte("d")},         // new
	}
	stat, err := w.Write("topic", "prov", pages2)
	if err != nil {
		t.Fatal(err)
	}
	if stat.Skipped != 1 {
		t.Errorf("Skipped = %d, want 1 (a.md)", stat.Skipped)
	}
	if stat.Written != 2 {
		t.Errorf("Written = %d, want 2 (b.md + d.md)", stat.Written)
	}
	if stat.Total != 3 {
		t.Errorf("Total = %d, want 3", stat.Total)
	}
}

func TestWrite_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	w := NewWriter(dir)

	pages := []provider.Page{
		{Path: "empty.md", Content: []byte{}},
		{Path: "real.md", Content: []byte("content")},
	}

	stat, err := w.Write("topic", "prov", pages)
	if err != nil {
		t.Fatal(err)
	}
	// Empty content pages are skipped entirely (not counted).
	if stat.Total != 1 {
		t.Errorf("Total = %d, want 1 (empty page not counted)", stat.Total)
	}
}

// findDateDir returns the first date-formatted subdirectory.
func findDateDir(t *testing.T, provDir string) string {
	t.Helper()
	entries, err := os.ReadDir(provDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() && e.Name() != "latest" {
			return filepath.Join(provDir, e.Name())
		}
	}
	t.Fatal("no date directory found")
	return ""
}
