package main

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
)

func main() {
	client := graphql.NewClient("https://api.github.com/graphql")
	q := `
		query {
			viewer {
				login
			}
		}`
	request := graphql.NewRequest(q)
	var response interface{}
	err := client.Run(context.Background(), request, &response)
	if err != nil {
		fmt.Println("problem", err)
	}
	fmt.Println(response)
}
