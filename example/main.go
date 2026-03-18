// Package main provides a small CLI example that connects to the
// Archivematica Storage Service and prints configured locations and
// pipelines.
//
// The example is intended to demonstrate both levels of this SDK:
// the higher-level ssclient wrapper API and the lower-level generated Kiota
// client exposed as an escape hatch. It lists locations and pipelines twice,
// first through the wrapper API and then through the raw Kiota request
// builders, so callers can compare both approaches side by side.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

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
	fs := flag.NewFlagSet("example", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		url  string
		user string
		key  string
	)
	fs.StringVar(&url, "url", "", "url, e.g. http://127.0.0.1:62081")
	fs.StringVar(&user, "user", "", "user, e.g.: test")
	fs.StringVar(&key, "key", "", "key, e.g.: test")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if url == "" || user == "" || key == "" {
		return flag.ErrHelp
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

	if err := app.wrapped(ctx); err != nil {
		return err
	}
	if err := app.rawClient(ctx); err != nil {
		return err
	}
	return nil
}

type application struct {
	client *ssclient.Client
	raw    *api.V2RequestBuilder
	stdout io.Writer
}

// wrapped prints the high-level client output.
func (app application) wrapped(ctx context.Context) error {
	app.printHeading("Using the wrapped client")

	listablePipelines, err := app.client.Pipelines().List(ctx, ssclient.ListPipelinesQuery{})
	if err != nil {
		return err
	}

	if err := app.printPipelines(listablePipelines); err != nil {
		return err
	}

	listable, err := app.client.Locations().List(ctx, ssclient.ListLocationsQuery{})
	if err != nil {
		return err
	}

	return app.printLocations(listable)
}

// rawClient prints the generated Kiota client output.
func (app application) rawClient(ctx context.Context) error {
	app.printHeading("Using the raw client")

	reqConfigPipeline := &api.V2PipelineEmptyPathSegmentRequestBuilderGetRequestConfiguration{}
	listablePipelines, err := app.raw.Pipeline().EmptyPathSegment().Get(ctx, reqConfigPipeline)
	if err != nil {
		return err
	}

	if err := app.printPipelines(listablePipelines); err != nil {
		return err
	}

	reqConfig := &api.V2LocationEmptyPathSegmentRequestBuilderGetRequestConfiguration{}
	listable, err := app.raw.Location().EmptyPathSegment().Get(ctx, reqConfig)
	if err != nil {
		return err
	}

	return app.printLocations(listable)
}

func (app application) printHeading(title string) {
	line := strings.Repeat("=", len(title))
	fmt.Fprintf(app.stdout, "\n%s\n%s\n%s\n\n", line, title, line)
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

// printPipelines consumes generated Kiota model types.
func (app application) printPipelines(listable models.PipelineListable) error {
	if listable == nil {
		fmt.Fprintf(app.stdout, "Found 0 pipelines!\n")
		return nil
	}

	count := len(listable.GetObjects())
	if meta := listable.GetMeta(); meta != nil && meta.GetTotalCount() != nil {
		count = int(*meta.GetTotalCount())
	}

	fmt.Fprintf(app.stdout, "Found %d pipelines!\n", count)
	for _, pipeline := range listable.GetObjects() {
		fmt.Fprintf(app.stdout, "» Pipeline %s with remote name %s.\n", valueOrEmpty(pipeline.GetUuid()), valueOrEmpty(pipeline.GetRemoteName()))
	}

	return nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
