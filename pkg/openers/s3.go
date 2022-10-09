//go:build !fileonly
// +build !fileonly

package openers

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/batmac/ccat/pkg/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	s3OpenerName        = "s3"
	s3OpenerDescription = "get an AWS s3 object via s3://"
)

type s3Opener struct {
	name, description string
}

func init() {
	register(&s3Opener{
		name:        s3OpenerName,
		description: s3OpenerDescription,
	})
}

func (f s3Opener) Name() string {
	return f.name
}

func (f s3Opener) Description() string {
	return f.description
}

func (f s3Opener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "s3://") {
		return 0.99
	}
	return 0
}

func (f s3Opener) Open(s string, _ bool) (io.ReadCloser, error) {
	ctx := context.Background()

	bucket, object := parseS3URI(s)
	log.Debugf("request to get %s in %s\n", object, bucket)

	if !isAWSEnvSet() {
		log.Debugf("  nothing found in env, setting a default profile and region\n")
		os.Setenv("AWS_PROFILE", "default")
		os.Setenv("AWS_REGION", "us-east-1")
	}

	// Load the Shared AWS Configuration (~/.aws/config)``
	log.Debugf(" LoadDefaultConfig...\n")
	// aws.LogRetries|aws.LogRequest|
	cfg, err := config.LoadDefaultConfig(ctx, config.WithClientLogMode(aws.LogDeprecatedUsage))
	if err != nil {
		return nil, err
	}

	if err := displayDebugInfoWithSTS(ctx, cfg); err != nil {
		return nil, err
	}

	log.Debugf(" creating S3 client...\n")
	s3Client := s3.NewFromConfig(cfg)
	log.Debugf("  GetObject...\n")
	o, err := s3Client.GetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(object),
		})
	if err != nil {
		return nil, err
	}

	return o.Body, nil
}

func parseS3URI(s string) (string, string) {
	s = strings.TrimPrefix(s, "s3://")
	pair := strings.SplitN(s, "/", 2)
	return pair[0], pair[1]
}

func displayDebugInfoWithSTS(ctx context.Context, cfg aws.Config) error {
	// Create an Amazon S3 service client
	log.Debugf(" creating STS client...\n")

	stsClient := sts.NewFromConfig(cfg) //nolint:contextcheck
	log.Debugf("  getting identity...\n")
	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	log.Debugf("  Account: %s, Arn: %s\n", aws.ToString(identity.Account), aws.ToString(identity.Arn))
	return nil
}

func isAWSEnvSet() bool {
	found := false

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "AWS_") {
			found = true
			log.Debugln("found in env: ", e)
		}
	}
	return found
}
