package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// Record sessions
	if err := s.RecordSession("alpha"); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordSession("beta"); err != nil {
		t.Fatal(err)
	}

	// List recent
	names, err := s.RecentSessions(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 recent sessions, got %d", len(names))
	}
	// beta should be more recent
	if names[0] != "beta" {
		t.Errorf("expected beta first (most recent), got %q", names[0])
	}

	// Re-record alpha to bump it
	if err := s.RecordSession("alpha"); err != nil {
		t.Fatal(err)
	}
	names, err = s.RecentSessions(10)
	if err != nil {
		t.Fatal(err)
	}
	if names[0] != "alpha" {
		t.Errorf("expected alpha first after re-record, got %q", names[0])
	}

	// Remove
	if err := s.RemoveSession("beta"); err != nil {
		t.Fatal(err)
	}
	names, err = s.RecentSessions(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 {
		t.Fatalf("expected 1 after remove, got %d", len(names))
	}
}

func TestStoreCorruptRecovery(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.db")

	// Write garbage to simulate a corrupt database.
	if err := os.WriteFile(path, []byte("not-a-sqlite-database"), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open should recover from corrupt db, got: %v", err)
	}
	defer s.Close()

	// Verify a .corrupt. backup was created.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range entries {
		if strings.Contains(e.Name(), ".corrupt.") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a .corrupt. backup file in the directory")
	}

	// The fresh database should be functional.
	if err := s.RecordSession("test"); err != nil {
		t.Fatalf("fresh db should work: %v", err)
	}
}

func TestStorePing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if err := s.Ping(); err != nil {
		t.Fatalf("Ping should succeed on a healthy store: %v", err)
	}
}

func TestStoreSchemaVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// Query schema_version directly.
	var version int
	err = s.db.QueryRow("SELECT version FROM schema_version").Scan(&version)
	if err != nil {
		t.Fatalf("schema_version table should exist: %v", err)
	}
	if version != 1 {
		t.Fatalf("expected schema version 1, got %d", version)
	}
}

// TestStoreSchemaVersionIdempotent verifies re-opening applies no duplicate migrations.
func TestStoreSchemaVersionIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	s1, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s1.Close()

	s2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	var version int
	if err := s2.db.QueryRow("SELECT version FROM schema_version").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if version != 1 {
		t.Fatalf("expected schema version 1 after re-open, got %d", version)
	}

	// Ensure only one row in schema_version.
	var count int
	if err := s2.db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row in schema_version, got %d", count)
	}
}
