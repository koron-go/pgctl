# `pg_ctl` command wrapper for golang

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/pgctl)](https://pkg.go.dev/github.com/koron-go/pgctl)
[![Actions/Go](https://github.com/koron-go/pgctl/workflows/Go/badge.svg)](https://github.com/koron-go/pgctl/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/pgctl)](https://goreportcard.com/report/github.com/koron-go/pgctl)

## Snippets

### Test with independent PostgreSQL

Use `tpg` sub package.

```golang
import (
    "database/sql"
    "github.com/koron-go/pgctl/tpg"
    _ "github.com/lib/pq"
    "testing"
)

func TestWithDB(t *testing.T) {
    // initialize database and start it
    srv := tpg.New(t)
    // remove all data when this test finished
    defer srv.Close()

    // connect to database
    db, err := sql.Open("postgres", srv.Name())
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // TODO: test with db!
}
```

### Basic usage

```golang
import (
    "database/sql"
    "github.com/koron-go/pgctl"
    _ "github.com/lib/pq"
)

// start PostgreSQL server
srv := pgctl.NewServer("dataDir")
if err := pgctl.Start(); err != nil {
    panic(err)
}
defer srv.Stop()

db, err := sql.Open("postgres", srv.Name())
if err != nil {
    panic(err)
}
defer db.Close()

// TODO: enjoy your work with "db"!
```

### Using raw API

```golang
import "github.com/koron-go/pgctl"

// initialize a database
err := pgctl.InitDB("dbdir", nil)

// start a database in background
err := pgctl.Start("dbdir", nil)

// check a database is running or not
err := pgctl.Status("dbdir")

// terminate a database
err := pgctl.Stop("dbdir")
```
