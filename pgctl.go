package pgctl

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

// InitDBOptions is a set of options of InitDB().
type InitDBOptions struct {
	User     string
	Encoding string
	Locale   string
}

// Options generates an string for "-o".
func (io *InitDBOptions) Options() string {
	args := make([]string, 0, 10)
	args = append(args,
		"-U", io.user(),
		"-A trust",
		"--encoding="+io.encoding(),
		"--locale="+io.locale())
	return strings.Join(args, " ")
}

func (io *InitDBOptions) user() string {
	if io.User == "" {
		return "postgres"
	}
	return io.User
}

func (io *InitDBOptions) encoding() string {
	if io.Encoding == "" {
		return "UTF8"
	}
	return io.Encoding
}

func (io *InitDBOptions) locale() string {
	if io.Locale == "" {
		return "C"
	}
	return io.Locale
}

// InitDB creates a dir and initiate as PostgreSQL database.
func InitDB(dir string, io *InitDBOptions) error {
	return initDBContext(context.Background(), dir, io)
}

// initDBContext creates a dir and initiate as PostgreSQL database with
// Context.
func initDBContext(ctx context.Context, dir string, io *InitDBOptions) error {
	if io == nil {
		io = &InitDBOptions{}
	}

	dd := filepath.Join(dir)
	_, err := os.Stat(dd)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return ErrAlreadyExists
	}

	cmd := exec.CommandContext(ctx, getPgCtl(), "initdb", "-s", "-D", dir, "-o", io.Options())
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// StartOptions is a set of options for Start().
type StartOptions struct {
	Port      uint16
	SocketDir string
	DBName    string
}

func (so *StartOptions) port() uint16 {
	if so.Port == 0 {
		return 5432
	}
	return so.Port
}

func (so *StartOptions) portString() string {
	return strconv.Itoa(int(so.port()))
}

func (so *StartOptions) host() string {
	return "127.0.0.1"
}

func (so *StartOptions) socketDir() string {
	if so.SocketDir == "" {
		return "''"
	}
	return so.SocketDir
}

// Options generates an string for "-o".
func (so *StartOptions) Options() string {
	args := make([]string, 0, 8)
	args = append(args, "-h", so.host(), "-F")
	if so.Port != 0 {
		args = append(args, "-p", so.portString())
	}
	args = append(args, "-k", so.socketDir())
	if so.DBName != "" {
		args = append(args, so.DBName)
	}
	return strings.Join(args, " ")
}

// Start starts PostgreSQL server on dir.
func Start(dir string, so *StartOptions) error {
	return startContext(context.Background(), dir, so)
}

// startContext starts PostgreSQL server on dir with Context.
func startContext(ctx context.Context, dir string, so *StartOptions) error {
	if so == nil {
		so = &StartOptions{}
	}

	err := statusContext(ctx, dir)
	if err == nil {
		return ErrAlreadyRunning
	}
	if err != ErrNotRunning {
		return err
	}

	cmd := exec.CommandContext(ctx, getPgCtl(), "start", "-s", "-D", dir, "-w", "-o", so.Options())
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
	cmd := exec.CommandContext(ctx, getPgCtl(), "status", "-D", dir)
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
	cmd := exec.CommandContext(ctx, getPgCtl(), "stop", "-s", "-D", dir)
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return ErrNotRunning
		}
		return err
	}
	return nil
}

// IsAvailable checks "pg_ctl" command is available or not.
func IsAvailable() bool {
	_, err := exec.LookPath(getPgCtl())
	return err == nil
}

func getPgCtl() string {
	dir, has := os.LookupEnv("POSTGRES_HOME")
	if !has {
		return "pg_ctl"
	}
	return filepath.Join(dir, "bin", "pg_ctl")
}

// Name returns data source name.
func Name(io *InitDBOptions, so *StartOptions) string {
	u := io.user()
	dbn := so.DBName
	if dbn == "" {
		dbn = u
	}
	return fmt.Sprintf("postgres://%[1]s@%[2]s:%[3]s/%[4]s", u, so.host(), so.portString(), dbn)
}
