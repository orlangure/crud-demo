// Package models includes functions to communicate with the storage layer.
package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Thing is a thing used in this application.
type Thing struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// String implements Stringer.
func (t *Thing) String() string {
	return fmt.Sprintf("%d: %s (%s)", t.ID, t.Name, t.Comment)
}

// Connect creates a new connection to the underlying storage layer.
func Connect(conn string) (*DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, fmt.Errorf("can't open sql connection: %w", err)
	}

	return &DB{db}, nil
}

// DB wraps a database connection and exposes allowed operations.
type DB struct {
	db *sql.DB
}

// CreateThing creates a new thing with the provided name and comment.
func (db *DB) CreateThing(name, comment string) error {
	query := `
		insert into things (name, comment)
		values (?, ?)
	`

	_, err := db.db.Exec(query, name, comment)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

// GetThingByName returns a thing with the provided name if exists.
func (db *DB) GetThingByName(name string) (*Thing, error) {
	query := `
		select id, name, comment
		from things
		where name = ?
		limit 1
	`
	row := db.db.QueryRow(query, name)

	t := &Thing{}

	err := row.Scan(&t.ID, &t.Name, &t.Comment)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return t, nil
}

// GetThingByID returns a thing with the provided id if exists.
func (db *DB) GetThingByID(id int) (*Thing, error) {
	query := `
		select id, name, comment
		from things
		where id = ?
		limit 1
	`
	row := db.db.QueryRow(query, id)

	t := &Thing{}

	err := row.Scan(&t.ID, &t.Name, &t.Comment)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return t, nil
}
