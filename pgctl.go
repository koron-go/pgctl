package pgctl

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	// ErrAlreadyExists means database is existing already
	ErrAlreadyExists = errors.New("database exists already")

	// ErrNotRunning means database is not running.
	ErrNotRunning = errors.New("database is not running")

	// ErrAlreadyRunning means database is running already.
	ErrAlreadyRunning = errors.New("database is running already")

	// ErrStartDatabase means failed to start database.
	ErrStartDatabase = errors.New("failed to start database")
)

// InitDBArgs is default arguments for InitDB()
var InitDBArgs = `-U postgres -A trust --encoding=UTF8 --locale=C`

// InitDB creates a dir and initiate as PostgreSQL database.
func InitDB(dir string) error {
	return initDBContext(context.Background(), dir)
}

// initDBContext creates a dir and initiate as PostgreSQL database with
// Context.
func initDBContext(ctx context.Context, dir string) error {
	dd := filepath.Join(dir)
	_, err := os.Stat(dd)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return ErrAlreadyExists
	}

	cmd := exec.CommandContext(ctx, "pg_ctl",
		"initdb", "-s", "-D", dir, "-o", InitDBArgs)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Start starts PostgreSQL server on dir.
func Start(dir string) error {
	return startContext(context.Background(), dir)
}

// startContext starts PostgreSQL server on dir with Context.
func startContext(ctx context.Context, dir string) error {
	err := statusContext(ctx, dir)
	if err == nil {
		return ErrAlreadyRunning
	}
	if err != ErrNotRunning {
		return err
	}

	cmd := exec.CommandContext(ctx, "pg_ctl", "start", "-s", "-D", dir)
	err = cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		return ErrStartDatabase
	}
	return err
}

// Status checks PostgreSQL server is running or not.
func Status(dir string) error {
	return statusContext(context.Background(), dir)
}

// statusContext checks PostgreSQL server is running or not, with Context.
func statusContext(ctx context.Context, dir string) error {
	cmd := exec.CommandContext(ctx, "pg_ctl", "status", "-D", dir)
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return ErrNotRunning
		}
		return err
	}
	return nil
}

// Stop stops PostgreSQL server on dir.
func Stop(dir string) error {
	return stopContext(context.Background(), dir)
}

// stopContext stops PostgreSQL server on dir with Context.
func stopContext(ctx context.Context, dir string) error {
	cmd := exec.CommandContext(ctx, "pg_ctl", "stop", "-s", "-D", dir)
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return ErrNotRunning
		}
		return err
	}
	return nil
}
