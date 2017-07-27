# dbstructload

This is a simple demonstration of an idea for simplifying the loading of data from database queries into structs. Struct fields that represent rows in the database are given `queryField` tags. Queries fields are then renamed to match these tags. Using the `reflect` package, query fields can then be matched to struct fields to simplify struct loading.

### Example

```go
import (
	"database/sql"
	"fmt"

	"github.com/mastahyeti/dbstructload"
)

// User represents a row from the `users` table.
type User struct {
	ID    uint64 `queryField:"User_id"`
	Login string `queryField:"User_login"`
}

func main() {
  db, err := sql.Open("mysql", "someaddr")
  if err != nil {
    panic(err)
  }

  const query = `
    SELECT
      users.id    AS User_id,
      users.login AS User_login,
    FROM users
    WHERE id=1
    LIMIT 1;
  `

  rows := dbstructload.Query(db, query)
  defer rows.Close()

  if ok := rows.Next(); !ok {
    panic("no row returned from query")
  }

  user := User{}
  if err := rows.Scan(&user); err != nil {
    panic(err)
  }

  fmt.Println("User 1 is ", user.Login)
}
```

A slightly more complex example can be found in [`/example/main.go`](example/main.go).
