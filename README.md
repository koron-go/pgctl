# `pg_ctl` command wrapper for golang

[![GoDoc](https://godoc.org/github.com/koron-go/pgctl?status.svg)](https://godoc.org/github.com/koron-go/pgctl)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/pgctl/master.svg)](https://circleci.com/gh/koron-go/pgctl/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/pgctl)](https://goreportcard.com/report/github.com/koron-go/pgctl)

## Snippets

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
err := pgctl.InitDB("dbdir")

// start a database in background
err := pgctl.Start("dbdir")

// check a database is running or not
err := pgctl.Status("dbdir")

// terminate a database
err := pgctl.Stop("dbdir")
```

## Tips

### Start PostgreSQL server on Debian or Ubuntu

It is required that write permission of /var/run/postgesql directory to start
PostgreSQL server on Debian or Ubuntu.

To do that, try these commands.

```console
$ sudo chmod o+w /var/run/postgresql
```

Or add a user to "postgres" group and logout then login.

```console
$ sudo adduser $USER postgres
```
