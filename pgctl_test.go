package pgctl

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestPgctl(t *testing.T) {
	_, err := exec.LookPath("pg_ctl")
	if err != nil {
		t.Skip("can't find pg_ctl", err)
	}

	tmp, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	db := filepath.Join(tmp, "db")

	// before InitDB
	err = Start(db)
	if err != ErrStartDatabase {
		t.Error("unexpected Start() result", err)
	}
	err = Stop(db)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
	err = Status(db)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}

	err = InitDB(db)
	if err != nil {
		t.Fatal("InitDB() failed", err)
	}

	// after InitDB() and before Start()
	err = InitDB(db)
	if err != ErrAlreadyExists {
		t.Error("InitDB() failed", err)
	}
	err = Stop(db)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
	err = Status(db)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}

	err = Start(db)
	if err != nil {
		t.Fatal("Start() failed", err)
	}
	time.Sleep(3 * time.Second)

	// after Start()
	err = Start(db)
	if err != ErrAlreadyRunning {
		t.Error("Start() failed", err)
	}
	err = Status(db)
	if err != nil {
		t.Error("Status() failed", err)
	}

	err = Stop(db)
	if err != nil {
		t.Fatal("Stop() failed", err)
	}

	// after Stop()
	err = Stop(db)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
	err = Status(db)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
}
