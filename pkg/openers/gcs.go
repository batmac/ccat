//go:build !fileonly
// +build !fileonly

package openers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/utils"
)

var (
	gcsOpenerName        = "gcs"
	gcsOpenerDescription = "get a GCP Cloud Storage object via gs://"
)

type gcsOpener struct {
	name, description string
}

func init() {
	register(&gcsOpener{
		name:        gcsOpenerName,
		description: gcsOpenerDescription,
	})
}

func (f gcsOpener) Name() string {
	return f.name
}

func (f gcsOpener) Description() string {
	return f.description
}

func (f gcsOpener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "gs://") {
		return 0.99
	}
	return 0
}

func (f gcsOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	ctx := context.Background()

	bucket, object := parseGcsURI(s)
	log.Debugf("request to get %s in %s\n", object, bucket)

	log.Debugf(" creating client...\n")

	// ~/.config/gcloud/application_default_credentials.json
	// or GOOGLE_APPLICATION_CREDENTIALS
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	// defer client.Close()

	/* 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	 */ // defer cancel()

	log.Debugf("  Get Object...\n")

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	// defer rc.Close()

	return utils.NewReadCloser(rc, func() error {
		log.Debugf("closing client...\n")
		return client.Close()
	}), nil
}

func parseGcsURI(s string) (string, string) {
	s = strings.TrimPrefix(s, "gs://")
	pair := strings.SplitN(s, "/", 2)
	return pair[0], pair[1]
}
