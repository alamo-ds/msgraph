package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/s-hammon/msgraph/graph"
)

func main() {
	groupId := "3cdbc4af-efe7-4f9f-9c2e-dcb8eb6128a9"
	ctx := context.Background()

	client := graph.NewClient()
	group, err := client.Groups().ById(groupId).Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		log.Fatalln("json.MarshalIndent:", err)
	}
	fmt.Println(string(data))
}
