package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

// Birthday represents a birthday record in the database
type Birthday struct {
	ID        int
	Name      string
	Month     int
	Day       int
	Gender    *string // Nullable for pronoun reference
	DiscordID *string // Nullable Discord user ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and initializes the schema
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Initialize schema
	if _, err := conn.Exec(schema); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// AddBirthday adds a new birthday to the database
func (db *DB) AddBirthday(name string, month, day int, gender, discordID *string) error {
	query := `INSERT INTO birthdays (name, month, day, gender, discord_id) VALUES (?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, name, month, day, gender, discordID)
	if err != nil {
		return fmt.Errorf("failed to add birthday: %w", err)
	}
	return nil
}

// GetBirthday gets a birthday by name
func (db *DB) GetBirthday(name string) (*Birthday, error) {
	query := `SELECT id, name, month, day, gender, discord_id, created_at, updated_at
	          FROM birthdays WHERE name = ?`

	var b Birthday
	err := db.conn.QueryRow(query, name).Scan(
		&b.ID, &b.Name, &b.Month, &b.Day, &b.Gender, &b.DiscordID, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get birthday: %w", err)
	}
	return &b, nil
}

// GetAllBirthdays returns all birthdays
func (db *DB) GetAllBirthdays() ([]Birthday, error) {
	query := `SELECT id, name, month, day, gender, discord_id, created_at, updated_at
	          FROM birthdays ORDER BY month, day`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query birthdays: %w", err)
	}
	defer rows.Close()

	var birthdays []Birthday
	for rows.Next() {
		var b Birthday
		if err := rows.Scan(&b.ID, &b.Name, &b.Month, &b.Day, &b.Gender, &b.DiscordID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan birthday: %w", err)
		}
		birthdays = append(birthdays, b)
	}

	return birthdays, nil
}

// GetBirthdaysByMonth returns all birthdays in a specific month
func (db *DB) GetBirthdaysByMonth(month int) ([]Birthday, error) {
	query := `SELECT id, name, month, day, gender, discord_id, created_at, updated_at
	          FROM birthdays WHERE month = ? ORDER BY day`

	rows, err := db.conn.Query(query, month)
	if err != nil {
		return nil, fmt.Errorf("failed to query birthdays by month: %w", err)
	}
	defer rows.Close()

	var birthdays []Birthday
	for rows.Next() {
		var b Birthday
		if err := rows.Scan(&b.ID, &b.Name, &b.Month, &b.Day, &b.Gender, &b.DiscordID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan birthday: %w", err)
		}
		birthdays = append(birthdays, b)
	}

	return birthdays, nil
}

// GetBirthdaysByDate returns all birthdays on a specific date
func (db *DB) GetBirthdaysByDate(month, day int) ([]Birthday, error) {
	query := `SELECT id, name, month, day, gender, discord_id, created_at, updated_at
	          FROM birthdays WHERE month = ? AND day = ?`

	rows, err := db.conn.Query(query, month, day)
	if err != nil {
		return nil, fmt.Errorf("failed to query birthdays by date: %w", err)
	}
	defer rows.Close()

	var birthdays []Birthday
	for rows.Next() {
		var b Birthday
		if err := rows.Scan(&b.ID, &b.Name, &b.Month, &b.Day, &b.Gender, &b.DiscordID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan birthday: %w", err)
		}
		birthdays = append(birthdays, b)
	}

	return birthdays, nil
}

// UpdateBirthday updates an existing birthday
func (db *DB) UpdateBirthday(name string, month, day int, gender, discordID *string) error {
	query := `UPDATE birthdays SET month = ?, day = ?, gender = ?, discord_id = ? WHERE name = ?`
	result, err := db.conn.Exec(query, month, day, gender, discordID, name)
	if err != nil {
		return fmt.Errorf("failed to update birthday: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no birthday found for %s", name)
	}

	return nil
}

// DeleteBirthday removes a birthday from the database
func (db *DB) DeleteBirthday(name string) error {
	query := `DELETE FROM birthdays WHERE name = ?`
	result, err := db.conn.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete birthday: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no birthday found for %s", name)
	}

	return nil
}

// GetPronoun returns the appropriate pronoun based on gender
func (b *Birthday) GetPronoun(subjectForm bool) string {
	// Default is they/them
	if b.Gender == nil {
		if subjectForm {
			return "they"
		}
		return "their"
	}

	switch *b.Gender {
	case "male":
		if subjectForm {
			return "he"
		}
		return "his"
	case "female":
		if subjectForm {
			return "she"
		}
		return "her"
	case "nonbinary", "other":
		if subjectForm {
			return "they"
		}
		return "their"
	default:
		if subjectForm {
			return "they"
		}
		return "their"
	}
}
