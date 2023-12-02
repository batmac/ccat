//go:build !fileonly
// +build !fileonly

package openers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/stringutils"

	"braces.dev/errtrace"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	mcOpenerName        = "mc"
	mcOpenerDescription = "get a Minio-compatible object via mc:// (use ~/.mc/config.json or env for credentials)"

	mcConfigPath = "/.mc/config.json"
	EnvVar       = Config{
		Endpoint:        "AWS_ENDPOINT",
		AccessKeyID:     "AWS_ACCESS_KEY_ID",
		SecretAccessKey: "AWS_SECRET_ACCESS_KEY",
	}
	lastResortEndpoint = "s3.amazonaws.com"
)

type Config struct {
	// mc stores scheme://endpoint as "url"
	Endpoint        string `json:"url"`
	AccessKeyID     string `json:"accessKey"`
	SecretAccessKey string `json:"secretKey"`
}

type mcOpener struct {
	name, description string
}

func init() {
	register(&mcOpener{
		name:        mcOpenerName,
		description: mcOpenerDescription,
	})
}

func (f mcOpener) Name() string {
	return f.name
}

func (f mcOpener) Description() string {
	return f.description
}

func (f mcOpener) Evaluate(s string) float32 {
	if strings.HasPrefix(s, "mc://") {
		return 0.99
	}
	return 0
}

func (f mcOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	ctx := context.Background()

	alias, bucket, object := parseMcURI(s)
	log.Debugf("request to get '%s' in '%s' from alias '%s'\n", object, bucket, alias)

	log.Debugf(" creating client...\n")

	c := getConfig(alias)
	log.Debugf("config: %+v\n", c)
	useSSL := !globalctx.GetBool("insecure")

	if c == (Config{}) {
		// print Help
		fmt.Printf("Could not find your Minio config,\n")
		fmt.Printf(" neither in %s\n", mcConfigPath)
		fmt.Printf(" neither in MC_HOST_<alias>\n")
		fmt.Printf(" neither in your env (%+v)", EnvVar)
		os.Exit(97)
	}

	if len(c.Endpoint) == 0 {
		c.Endpoint = lastResortEndpoint
	}

	minioClient, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	log.Debugf("  Get Object %s@%s...\n", object, bucket)

	o, err := minioClient.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	return o, nil
}

func parseMcURI(s string) (string, string, string) {
	s = strings.TrimPrefix(s, "mc://")
	pair := strings.SplitN(s, "/", 3)
	// alias, bucket, object
	return pair[0], pair[1], pair[2]
}

func getConfig(alias string) Config {
	if len(alias) == 0 {
		log.Debugln("no alias given, looking in env...")
		return getConfigFromFallback()
	}
	confFile := getConfigFromMCFile(alias)
	confEnv := getConfigFromMCEnv(alias)
	return mergeConfig(confFile, confEnv)
}

func mergeConfig(c, c2 Config) Config {
	c1 := c
	if len(c2.Endpoint) > 0 {
		c1.Endpoint = c2.Endpoint
	}
	if len(c2.AccessKeyID) > 0 {
		c1.AccessKeyID = c2.AccessKeyID
	}
	if len(c2.SecretAccessKey) > 0 {
		c1.SecretAccessKey = c2.SecretAccessKey
	}
	// log.Debugf("merged %+v\n", c1)
	return c1
}

func getConfigFromMCFile(alias string) Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return Config{}
	}

	type ConfFileData struct {
		Aliases map[string]Config `json:"aliases"`
		Version string
	}

	confJ, err := os.ReadFile(filepath.Clean(home + mcConfigPath))
	if err != nil {
		log.Debugln(err)
		return Config{}
	}
	var conf ConfFileData
	err = json.Unmarshal(confJ, &conf)
	if err != nil {
		log.Debugln(err)
		return Config{}
	}

	if a, ok := conf.Aliases[alias]; ok {
		log.Debugf("found a config for '%s' in mc config (version %s)\n", alias, conf.Version)
		return Config{stringutils.RemoveScheme(a.Endpoint), a.AccessKeyID, a.SecretAccessKey}
	}

	return Config{}
}

func getConfigFromMCEnv(alias string) Config {
	var c Config
	host := "MC_HOST_" + alias
	if val, ok := os.LookupEnv(host); ok {
		log.Debugf("config found in %s\n", host)
		u, err := url.Parse(val)
		if err == nil {
			if len(u.Host) != 0 {
				c.Endpoint = u.Host
			}
			if username := u.User.Username(); len(username) > 0 {
				c.AccessKeyID = username
			}
			if pass, b := u.User.Password(); b {
				c.SecretAccessKey = pass
			} else {
				log.Printf("failed to parse %s\n", host)
			}
		}
	}
	// log.Debugf("config found: %+v\n", c)
	return c
}

func getConfigFromFallback() Config {
	return Config{
		Endpoint:        os.Getenv(EnvVar.Endpoint),
		AccessKeyID:     os.Getenv(EnvVar.AccessKeyID),
		SecretAccessKey: os.Getenv(EnvVar.SecretAccessKey),
	}
}
