# ent

## Setup A GoEh Environment

```bash
❯ go version
go version go1.16.6 darwin/amd64

❯ go mod init github.com/k3forx/ent
```

## Installation

```bash
❯ go get entgo.io/ent/cmd/ent
```

## Create Your First Schema

```bash
❯ go run entgo.io/ent/cmd/ent init User
```

The above command creates the following file.

```go
package schema

import "entgo.io/ent"

// User holds the schema definition for the User entity.
type User struct {
        ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
        return nil
}

// Edges of the User.
func (User) Edges() []ent.Edge {
        return nil
}
```

Add 2 fields to the `User` schema:

```go
package schema

// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Positive(),
        field.String("name").
            Default("unknown"),
    }
}
```

Run `go generate` from the root directory of the project as follows:

```bash
❯ go generate ./ent
```

This produces the following files:

```bash
❯ tree ./ent
./ent
├── client.go
├── config.go
├── context.go
│   └── user.go
... truncated
├── tx.go
├── user
│   ├── user.go
│   └── where.go
├── user.go
├── user_create.go
├── user_delete.go
├── user_query.go
└── user_update.go

7 directories, 22 files
```

## Create Your First Entity

To get started, create a new `ent.Client`. For this example, we will use MySQL.

```bash
❯ cat main.go
package main

import (
        "context"
        "fmt"
        "log"

        "github.com/k3forx/ent/ent"

        _ "github.com/go-sql-driver/mysql"
)

func main() {
        client, err := ent.Open("mysql", "mysql:mysql@tcp(localhost:3306)/test?parseTime=True")
        if err != nil {
                log.Fatalf("failed opening connection to mysql: %v", err)
        }
        defer client.Close()

        // Run the auto migration tool
        if err := client.Schema.Create(context.Background()); err != nil {
                log.Fatalf("failed creating schema resources: %v", err)
        }
}
```

Now, we're ready to create our user. Let's call this function `CreateUser` for the sake of example:

```bash
func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
        u, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)
        if err != nil {
                return nil, fmt.Errorf("failed creating user: %w", err)
        }
        log.Println("user was created: ", u)
        return u, nil
}
```

## Run the script

```bash
❯ go run ./main.go
2021/08/11 23:36:19 user was created:  User(id=1, age=30, name=a8m)
user: User(id=1, age=30, name=a8m)%
```
