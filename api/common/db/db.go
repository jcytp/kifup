package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jcytp/kifup-api/common/env"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

var db *DB

func CheckDBFileExists() bool {
	_, err := os.Stat(env.DatabasePath())
	return !os.IsNotExist(err)
}

func New() {
	sqlt, err := sql.Open("sqlite", env.DatabasePath())
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqlt.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Error enabling foreign key constraints:", err)
	}

	db = &DB{sqlt}
}

func CheckConnection() {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	slog.Info("Successfully connected to the database")
}

func Close() {
	db.DB.Close()
}

// ------------------------------------------------------------

func Exec(sql string, args ...any) (sql.Result, error) {
	result, err := db.DB.Exec(sql, args...)
	if err != nil {
		slog.Error("SQL Error", "query", formatSQL(sql), "params", formatParams(args), "error", err)
	}
	return result, err
}

func QueryRow(sql string, args ...any) *sql.Row {
	return db.DB.QueryRow(sql, args...)
}

func Query(sql string, args ...any) (*sql.Rows, error) {
	rows, err := db.DB.Query(sql, args...)
	if err != nil {
		slog.Error("SQL Error", "query", formatSQL(sql), "params", formatParams(args), "error", err)
	}
	return rows, err
}

func formatSQL(sql string) string {
	sql = strings.ReplaceAll(sql, "\t", " ")
	lines := strings.Split(sql, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, " ")
}

func formatParams(args []any) string {
	params := make([]string, len(args))
	for i, arg := range args {
		params[i] = fmt.Sprintf("%v", arg)
	}
	return "[" + strings.Join(params, ", ") + "]"
}

func CheckAffectedRows(result sql.Result, cnt int64) error {
	cntAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if cntAffected != cnt {
		return fmt.Errorf("affected rows error")
	}
	return nil
}

// ------------------------------------------------------------

func PInt64(n int64) *int64 {
	return &n
}

func PTime(t time.Time) *time.Time {
	return &t
}
