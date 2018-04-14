# `pg_ctl` command wrapper for golang

[![GoDoc](https://godoc.org/github.com/koron-go/pgctl?status.svg)](https://godoc.org/github.com/koron-go/pgctl)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/pgctl.svg)](https://circleci.com/gh/koron-go/pgctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/pgctl)](https://goreportcard.com/report/github.com/koron-go/pgctl)

## Snippets

```golang
import "github.com/koron-go/pgctl"

// initialize a database
err := pgctl.InitDB("dbdir")

// start a database in background
err := pgctl.Start("dbdir")

// check a database is running or not
err := pgctl.Status("dbdir")

// terminate a database
err := pgctl.Stop("dbdri")
```
