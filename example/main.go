package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.artefactual.dev/ssclient"
)

func main() {
	ctx := context.Background()

	// Create client.
	client, err := ssclient.New(&http.Client{}, "http://127.0.0.1:62081", "test", "test")
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	// List locations.
	listable, err := client.Api().V2().Location().Get(ctx, nil)
	if err != nil {
		log.Fatalf("Error listing locations: %v\n", err)
	}

	// Print locations.
	fmt.Printf("Found %d locations!\n", *listable.GetMeta().GetTotalCount())
	for _, location := range listable.GetObjects() {
		fmt.Printf("» Location %s with purpose %s.\n", *location.GetUuid(), location.GetPurpose())
	}
}
