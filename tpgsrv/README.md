# tpgsrv - PostgreSQL server for test

`tpgsrv` helps to create independent PostgreSQL server instance for test.

## Example

```golang
import (
    "database/sql"
    "github.com/koron-go/pgctl/tpgsrv"
    _ "github.com/lib/pq"
    "testing"
)

func TestWithDB(t *testing.T) {
    // initialize database and start it
    srv := tpgsrv.New(t)
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
