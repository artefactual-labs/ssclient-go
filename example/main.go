// Package main provides a small CLI example that connects to the
// Archivematica Storage Service and prints the configured locations.
//
// The example is intended to demonstrate both levels of this SDK:
// the higher-level ssclient wrapper API and the lower-level generated Kiota
// client exposed as an escape hatch. It lists the same locations twice, first
// through the wrapper API and then through the raw Kiota request builders, so
// callers can compare both approaches side by side.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/api"
	"go.artefactual.dev/ssclient/kiota/models"
)

const usage = "Usage: example -url=http://127.0.0.1:62081 -user=test -key=test"

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args[1:]); err == flag.ErrHelp {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		os.Exit(2)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, stdout io.Writer, args []string) error {
	if len(args) < 2 {
		return flag.ErrHelp
	}

	var (
		url  string
		user string
		key  string
	)
	flag.StringVar(&url, "url", "", "url, e.g. http://127.0.0.1:62081")
	flag.StringVar(&user, "user", "", "user, e.g.: test")
	flag.StringVar(&key, "key", "", "key, e.g.: test")
	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	client, err := ssclient.New(ssclient.Config{
		BaseURL:  url,
		Username: user,
		Key:      key,
	})
	if err != nil {
		return fmt.Errorf("create ssclient: %v", err)
	}

	app := application{
		client: client,
		raw:    client.Raw().Api().V2(),
		stdout: stdout,
	}

	if err := app.locations(ctx); err != nil {
		return err
	}

	return app.locationsRaw(ctx)
}

type application struct {
	client *ssclient.Client
	raw    *api.V2RequestBuilder
	stdout io.Writer
}

// locations prints a list of locations found in the storage server.
func (app application) locations(ctx context.Context) error {
	listable, err := app.client.Locations().List(ctx, ssclient.ListLocationsQuery{})
	if err != nil {
		return err
	}

	return app.printLocations(listable)
}

// locationsRaw prints a list of locations using the generated Kiota client.
func (app application) locationsRaw(ctx context.Context) error {
	reqConfig := &api.V2LocationEmptyPathSegmentRequestBuilderGetRequestConfiguration{}
	listable, err := app.raw.Location().EmptyPathSegment().Get(ctx, reqConfig)
	if err != nil {
		return err
	}

	return app.printLocations(listable)
}

// printLocations consumes generated Kiota model types. The example keeps
// models shared across the wrapper layer and the raw client to show the
// current public API compromise: higher-level operations, but common schema
// types.
func (app application) printLocations(listable models.LocationListable) error {
	if listable == nil {
		fmt.Fprintf(app.stdout, "Found 0 locations!\n")
		return nil
	}

	count := len(listable.GetObjects())
	if meta := listable.GetMeta(); meta != nil && meta.GetTotalCount() != nil {
		count = int(*meta.GetTotalCount())
	}

	fmt.Fprintf(app.stdout, "Found %d locations!\n", count)
	for _, location := range listable.GetObjects() {
		fmt.Fprintf(app.stdout, "» Location %s with purpose %s.\n", *location.GetUuid(), location.GetPurpose())
	}

	return nil
}
