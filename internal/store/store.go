package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store manages local SQLite state for agentmux.
type Store struct {
	db *sql.DB
}

// DefaultPath returns the default database path.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "agentmux", "state.db")
}

// Open opens or creates the SQLite database. If the existing database is
// corrupt (migrate fails), the file is renamed to state.db.corrupt.{timestamp}
// and a fresh database is created.
func Open(path string) (*Store, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("store: mkdir: %w", err)
	}

	s, err := openAndMigrate(path)
	if err != nil {
		// Attempt corrupt-database recovery if the file exists.
		if _, statErr := os.Stat(path); statErr == nil {
			ts := time.Now().UTC().Format("20060102T150405Z")
			corrupt := path + ".corrupt." + ts
			if renameErr := os.Rename(path, corrupt); renameErr != nil {
				return nil, fmt.Errorf("store: rename corrupt db: %w (original error: %w)", renameErr, err)
			}
			fmt.Fprintf(os.Stderr, "store: corrupt database moved to %s; creating fresh database\n", corrupt)

			s, err = openAndMigrate(path)
			if err != nil {
				return nil, fmt.Errorf("store: open after recovery: %w", err)
			}
			return s, nil
		}
		return nil, err
	}
	return s, nil
}

func openAndMigrate(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("store: open: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("store: migrate: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}

// RecordSession records a session access event.
func (s *Store) RecordSession(name string) error {
	_, err := s.db.Exec(`
		INSERT INTO recent_sessions (name, last_used)
		VALUES (?, ?)
		ON CONFLICT(name) DO UPDATE SET last_used = excluded.last_used
	`, name, time.Now().UnixNano())
	return err
}

// RecentSessions returns the most recently used session names.
func (s *Store) RecentSessions(limit int) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT name FROM recent_sessions
		ORDER BY last_used DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// RemoveSession removes a session from recents.
func (s *Store) RemoveSession(name string) error {
	_, err := s.db.Exec(`DELETE FROM recent_sessions WHERE name = ?`, name)
	return err
}

// Ping verifies the database connection is alive.
func (s *Store) Ping() error {
	var n int
	return s.db.QueryRow("SELECT 1").Scan(&n)
}

func migrate(db *sql.DB) error {
	// Ensure schema_version table exists.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER NOT NULL
		)
	`); err != nil {
		return err
	}

	// Seed version row if missing.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_version").Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := db.Exec("INSERT INTO schema_version (version) VALUES (0)"); err != nil {
			return err
		}
	}

	var ver int
	if err := db.QueryRow("SELECT version FROM schema_version").Scan(&ver); err != nil {
		return err
	}

	// Apply migrations incrementally.
	migrations := []func(*sql.DB) error{
		// Version 0 -> 1: create recent_sessions table.
		func(d *sql.DB) error {
			_, err := d.Exec(`
				CREATE TABLE IF NOT EXISTS recent_sessions (
					name      TEXT PRIMARY KEY,
					last_used INTEGER NOT NULL
				)
			`)
			return err
		},
	}

	for i := ver; i < len(migrations); i++ {
		if err := migrations[i](db); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
		if _, err := db.Exec("UPDATE schema_version SET version = ?", i+1); err != nil {
			return fmt.Errorf("update schema_version to %d: %w", i+1, err)
		}
	}

	return nil
}
