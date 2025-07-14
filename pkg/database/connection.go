package database

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Connections holds all database connections
type Connections struct {
	Characters *sql.DB
	Auth       *sql.DB
	World      *sql.DB
}

// NewConnections creates new database connections for all three databases
func NewConnections(dsn string) (*Connections, error) {
	// Create connection to characters database
	charactersDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the characters connection
	if err := charactersDB.Ping(); err != nil {
		return nil, err
	}

	// Create connection to auth database by modifying the DSN
	authDSN := strings.Replace(dsn, "/acore_characters?", "/acore_auth?", 1)
	authDB, err := sql.Open("mysql", authDSN)
	if err != nil {
		closeWithLog(charactersDB, "characters")
		return nil, err
	}

	// Test the auth connection
	if err := authDB.Ping(); err != nil {
		closeWithLog(charactersDB, "characters")
		closeWithLog(authDB, "auth")
		return nil, err
	}

	// Create connection to world database by modifying the DSN
	worldDSN := strings.Replace(dsn, "/acore_characters?", "/acore_world?", 1)
	worldDB, err := sql.Open("mysql", worldDSN)
	if err != nil {
		closeWithLog(charactersDB, "characters")
		closeWithLog(authDB, "auth")
		return nil, err
	}

	// Test the world connection
	if err := worldDB.Ping(); err != nil {
		closeWithLog(charactersDB, "characters")
		closeWithLog(authDB, "auth")
		closeWithLog(worldDB, "world")
		return nil, err
	}

	return &Connections{
		Characters: charactersDB,
		Auth:       authDB,
		World:      worldDB,
	}, nil
}

// Close closes all database connections
func (c *Connections) Close() {
	if c.Characters != nil {
		closeWithLog(c.Characters, "characters")
	}
	if c.Auth != nil {
		closeWithLog(c.Auth, "auth")
	}
	if c.World != nil {
		closeWithLog(c.World, "world")
	}
}

// Helper functions for error handling
func closeWithLog(db *sql.DB, name string) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing %s database: %v", name, err)
	}
}

// CloseRowsWithLog closes rows with error logging
func CloseRowsWithLog(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("Error closing rows: %v", err)
	}
}
