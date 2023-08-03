package main

import (
	stdctx "context"
	"fmt"
	"web_dev/context"
	"web_dev/models"
)

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "test@main.com",
	}

	ctx = context.WithUser(ctx, &user)

	retrieveUser := context.User(ctx)

	fmt.Println(retrieveUser.Email)
}
