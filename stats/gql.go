package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/machinebox/graphql"
)

func main() {
	query := flag.String("q", "query { viewer { login } }", "GraphQL query string")
	token := flag.String("token", "<token>", "Github token")
	flag.Parse()

	client := graphql.NewClient("https://api.github.com/graphql")
	request := graphql.NewRequest(*query)
	request.Header.Set("Authorization", fmt.Sprintf("bearer %v", *token))
	var response interface{}
	err := client.Run(context.Background(), request, &response)
	if err != nil {
		fmt.Println("problem", err)
	}
	fmt.Println(response)
}
