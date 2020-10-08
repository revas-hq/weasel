package main

import (
	"context"
	"encoding/base64"
	"flag"
	"net/http"
	"os"
	"time"

	kit_log "github.com/go-kit/kit/log"

	"cloud.google.com/go/storage"
	"github.com/allegro/bigcache"
	"github.com/peterbourgon/ff/v3"
	"google.golang.org/api/option"

	weasel "github.com/revas-hq/weasel/pkg"
	weasel_inmem "github.com/revas-hq/weasel/pkg/service"
	weasel_gcloud "github.com/revas-hq/weasel/pkg/storage/gcloud"
	weasel_http "github.com/revas-hq/weasel/pkg/transport/http"
)

func main() {
	var logger kit_log.Logger
	{
		logger = kit_log.NewLogfmtLogger(kit_log.NewSyncWriter(os.Stderr))
		logger = kit_log.With(logger, "when", kit_log.DefaultTimestampUTC)
		logger = kit_log.With(logger, "where", kit_log.DefaultCaller)
		logger = kit_log.With(logger, "system", "weasel")
	}

	fs := flag.NewFlagSet("revas-weasel", flag.ExitOnError)
	var (
		GCloudCredentials64 = fs.String("gcloud-credentials", "", "Google Cloud IAM | base 64 json key with storage read access")
		GCloudBucket        = fs.String("gcloud-bucket", "", "Google Cloud Storage | bucket name")
	)
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("REVAS_WEASEL"))
	if err != nil {
		_ = logger.Log("err", err)
		panic("Unable to parse flags")
	}

	if *GCloudBucket == "" {
		panic("Google Cloud Storage bucket name cannot be empty")
	}

	GCloudCredentials, err := base64.StdEncoding.DecodeString(*GCloudCredentials64)
	if err != nil {
		_ = logger.Log("err", err)
		panic("Google Cloud data credentials should be base 64 encoded json files")
	}

	ctx := context.Background()

	var client *storage.Client
	{
		if len(GCloudCredentials) == 0 {
			_ = logger.Log("info", "using Google Cloud standard credentials")
			client, err = storage.NewClient(ctx)
		} else {
			_ = logger.Log("info", "using Google Cloud environment credentials")
			client, err = storage.NewClient(ctx, option.WithCredentialsJSON(GCloudCredentials))
		}
		if err != nil {
			_ = logger.Log("err", err)
			panic("Unable to connect to Google Cloud Storage services")
		}
	}

	bcache, err := bigcache.NewBigCache(bigcache.DefaultConfig(86400 * time.Second))
	if err != nil {
		_ = logger.Log("err", err)
		panic("Unable to create BigCache cache")
	}

	storage := weasel_gcloud.NewStorage(kit_log.With(logger, "subsystem", "storage"), client, *GCloudBucket)
	// TODO verify actual bucket read capability

	service := weasel.NewService(kit_log.With(logger, "subsystem", "service"), storage)
	service = weasel_inmem.NewCache(kit_log.With(logger, "subsystem", "cache"), bcache, service)

	endpoints := weasel.Endpoints{
		GetObject: weasel.NewGetObjectEndpoint(service),
	}

	handler := weasel_http.NewHTTPHandler(kit_log.With(logger, "subsystem", "http"), &endpoints)
	handler = weasel_http.NewCORSMiddleware()(handler)
	handler = weasel_http.NewForceTLSMiddleware()(handler)
	handler = weasel_http.NewForceSTSMiddleware()(handler)
	handler = weasel_http.NewCheckMethodMiddleware()(handler)
	http.Handle("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		_ = logger.Log("info", "defaulting to port "+port)
	}
	_ = logger.Log("info", "listening on port "+port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		_ = logger.Log("err", err)
	}
}
