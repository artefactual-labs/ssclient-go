package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota"
)

const usage = "Usage: example -url=http://127.0.0.1:62081 -user=test -key=test"

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, out io.Writer, args []string) error {
	if len(args) < 2 {
		return errors.New(usage)
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

	client, err := ssclient.New(&http.Client{}, url, user, key)
	if err != nil {
		return fmt.Errorf("create ssclient: %v", err)
	}

	app := application{client, out}

	return app.locations(ctx)
}

type application struct {
	client *kiota.Client
	stdout io.Writer
}

// locations prints a list of locations found in the storage server.
func (app application) locations(ctx context.Context) error {
	listable, err := app.client.Api().V2().Location().Get(ctx, nil)
	if err != nil {
		return err
	}

	fmt.Fprintf(app.stdout, "Found %d locations!\n", *listable.GetMeta().GetTotalCount())
	for _, location := range listable.GetObjects() {
		fmt.Fprintf(app.stdout, "» Location %s with purpose %s.\n", *location.GetUuid(), location.GetPurpose())
	}

	return nil
}
