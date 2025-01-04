package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jcytp/kifup-api/common/aws"
	"github.com/jcytp/kifup-api/common/env"
	_ "modernc.org/sqlite"
)

var db *sql.DB

const (
	currentDBKey       = "private/database/current.db"
	backupDBKeyPattern = "private/database/%s/%s.db" // YYYYMM/YYYYMMDD_HHMMSS.db
	backupCyclePeriod  = 4 * time.Hour
)

func CheckDBFileExists() bool {
	_, err := os.Stat(env.DatabasePath())
	return !os.IsNotExist(err)
}

func DownloadDB() error {
	if env.IsDevelopment() {
		return fmt.Errorf("development environment cannot access to S3")
	}

	if err := aws.S3DownloadToFile(env.S3BucketName(), currentDBKey, env.DatabasePath()); err != nil {
		slog.Warn("Failed to download database from S3", "error", err)
		return fmt.Errorf("failed to download database: %v", err)
	}

	slog.Info("Successfully download database from S3")
	return nil
}

func UploadDB() {
	if env.IsDevelopment() {
		return
	}

	now := time.Now()
	backupKey := fmt.Sprintf(backupDBKeyPattern, now.Format("200601"), now.Format("20060102_150405"))

	if err := aws.S3UploadFromFile(env.S3BucketName(), backupKey, env.DatabasePath()); err != nil {
		slog.Error("Failed to create database backup", "error", err)
		return
	}

	slog.Info("Successfully uploaded database to S3", "backup", backupKey)
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

	db = sqlt
}

func CheckConnection() {
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	slog.Info("Successfully connected to the database")
}

func Close() {
	db.Close()
}

func StartBackupCycle() {
	ticker := time.NewTicker(backupCyclePeriod)
	go func() {
		for range ticker.C {
			UploadDB()
		}
	}()
}

func ScheduleFinalBackup() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		slog.Info("Shutting down, uploading final DB backup...")
		UploadDB()
		os.Exit(0)
	}()
}

// ------------------------------------------------------------

func Exec(sql string, args ...any) (sql.Result, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		slog.Error("SQL Error", "query", formatSQL(sql), "params", formatParams(args), "error", err)
	}
	return result, err
}

func QueryRow(sql string, args ...any) *sql.Row {
	return db.QueryRow(sql, args...)
}

func Query(sql string, args ...any) (*sql.Rows, error) {
	rows, err := db.Query(sql, args...)
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
