# ent

## Setup A Go Environment

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

- `main.go`

```golang
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

```golang
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

## Add Your First Edge (Relation)

We want to declare an edge (relation) to another entity in the schema.
Let's create two additional entities named `Car` and `Group` with a few fields. We use `ent` CLI to generate this initial schemas:

```bash
❯ go run entgo.io/ent/cmd/ent init Car Group
```

And then we add the rest of the fields manually:

- `./ent/schema/car.go`

```golang
// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("model"),
		field.Time("registered_at"),
	}
}
```

- `ent/schema/group.go`

```golang
// Fields of the Group.
func (Group) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            // Regexp validation for group name.
            Match(regexp.MustCompile("[a-zA-Z_]+$")),
    }
}
```

Let's define our first relation. An edge from `User` to `Car` defining that a user can **have 1 or more** cars, but a car **has only one** owner (one-to-many relation).

Let's add the `"cars"` edge to the `User` schema, and run `go generation ./ent`:

- `ent/schema/user.go`

```golang
// Edges of the User.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("cars", Car.Type),
    }
}
```

We continue our example by creating 2 cars and adding them to a user.

- `main.go`

```golang
func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
    // Create a new car with model "Tesla".
    tesla, err := client.Car.
        Create().
        SetModel("Tesla").
        SetRegisteredAt(time.Now()).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating car: %w", err)
    }
    log.Println("car was created: ", tesla)

    // Create a new car with model "Ford".
    ford, err := client.Car.
        Create().
        SetModel("Ford").
        SetRegisteredAt(time.Now()).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating car: %w", err)
    }
    log.Println("car was created: ", ford)

    // Create a new user, and add it the 2 cars.
    a8m, err := client.User.
        Create().
        SetAge(30).
        SetName("a8m").
        AddCars(tesla, ford).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating user: %w", err)
    }
    log.Println("user was created: ", a8m)
    return a8m, nil
}
```

But what about querying the `cars` edge (relation)? Here's how we do it:

- `main.go`

```golang
func QueryCars(ctx context.Context, a8m *ent.User) error {
    cars, err := a8m.QueryCars().All(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    log.Println("returned cars:", cars)

    // What about filtering specific cars.
    ford, err := a8m.QueryCars().
        Where(car.Model("Ford")).
        Only(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    log.Println(ford)
    return nil
}
```

## Add Your First Inverse Edge (BackRef)

Assume we have a `Car` object and we want to get its owner; the user that this car belongs to. For this, we have another type of edge called "inverse edge" that is defined using the `edge.Form` function.

![image](https://user-images.githubusercontent.com/45956169/129462281-fc3470ec-5594-450a-9e34-690da6d4e7c2.png)

The new edge created in the diagram above is translucent, to emphasize that we don't create another edge in the database. It's just a back-reference to the real edge (relation).

Let's add an inverse edge named `owner` to the `Car` schema, reference it to the `cars` edge in the `User` schema, and run `go generate ./ent`.

- `ent/schema/car.go`

```golang
// Edges of the Car.
func (Car) Edges() []ent.Edge {
    return []ent.Edge{
        // Create an inverse-edge called "owner" of type `User`
        // and reference it to the "cars" edge (in User schema)
        // explicitly using the `Ref` method.
        edge.From("owner", User.Type).
            Ref("cars").
            // setting the edge to unique, ensure
            // that a car can have only one owner.
            Unique(),
    }
}
```

We'll continue the user/cars example above by querying the inverse edge.

- `main.go`

```golang
func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
    cars, err := a8m.QueryCars().All(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    // Query the inverse edge.
    for _, ca := range cars {
        owner, err := ca.QueryOwner().Only(ctx)
        if err != nil {
            return fmt.Errorf("failed querying car %q owner: %w", ca.Model, err)
        }
        log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
    }
    return nil
}
```
