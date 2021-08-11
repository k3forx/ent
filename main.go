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

	ctx := context.Background()

	u, err := CreateUser(ctx, client)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	fmt.Printf("user: %v", u)
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}
