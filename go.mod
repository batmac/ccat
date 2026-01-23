module github.com/batmac/ccat

go 1.25.0

require (
	cloud.google.com/go/storage v1.59.1
	git.sr.ht/~adnano/go-gemini v0.2.6
	github.com/JohannesKaufmann/html-to-markdown v1.6.0
	github.com/OneOfOne/xxhash v1.2.8
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/aixiansheng/lzfse v0.2.4
	github.com/alecthomas/chroma/v2 v2.23.0
	github.com/atotto/clipboard v0.1.4
	github.com/aws/aws-sdk-go-v2 v1.40.0
	github.com/aws/aws-sdk-go-v2/config v1.32.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.92.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.2
	github.com/batmac/go-curl v0.0.1
	github.com/cosnicolaou/pbzip2 v1.0.6
	github.com/creativeprojects/go-selfupdate v1.5.2
	github.com/docker/go-units v0.5.0
	github.com/dsnet/compress v0.0.1
	github.com/eliukblau/pixterm v1.3.2
	github.com/gabriel-vasile/mimetype v1.4.12
	github.com/gage-technologies/mistral-go v1.1.0
	github.com/gdamore/tcell/v2 v2.13.7
	github.com/google/generative-ai-go v0.20.1
	github.com/google/renameio/v2 v2.0.2
	github.com/gowebpki/jcs v1.0.1
	github.com/hbollon/go-edlib v1.7.0
	github.com/kevinburke/nacl v0.9.0
	github.com/klauspost/compress v1.18.3
	github.com/klauspost/pgzip v1.2.6
	github.com/kofalt/go-memoize v0.0.0-20240506050413-9e5eb99a0f2a
	github.com/magefile/mage v1.15.0
	github.com/minio/minio-go/v7 v7.0.98
	github.com/minio/minlz v1.0.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mmcdole/gofeed v1.3.0
	github.com/muesli/reflow v0.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pierrec/lz4/v4 v4.1.25
	github.com/psanford/wormhole-william v1.0.8
	github.com/rivo/tview v0.42.0
	github.com/robert-nix/ansihtml v1.0.1
	github.com/sashabaranov/go-openai v1.41.2
	github.com/spf13/pflag v1.0.10
	github.com/tetratelabs/wazero v1.11.0
	github.com/titanous/json5 v1.0.0
	github.com/tmc/keyring v0.0.0-20230418032330-0c8bdba76fa8
	github.com/traefik/yaegi v0.16.1
	github.com/ulikunitz/xz v0.5.15
	github.com/zeebo/xxh3 v1.0.2
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b
	golang.org/x/term v0.39.0
	google.golang.org/api v0.261.0
	howett.net/plist v1.0.1
	sigs.k8s.io/yaml v1.6.0
)

require (
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/ai v0.8.0 // indirect
	cloud.google.com/go/auth v0.18.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.5.3 // indirect
	cloud.google.com/go/longrunning v0.7.0 // indirect
	cloud.google.com/go/monitoring v1.24.3 // indirect
	code.gitea.io/sdk/gitea v0.22.1 // indirect
	github.com/42wim/httpsig v1.2.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.30.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.54.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.54.0 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/PaesslerAG/gval v1.0.0 // indirect
	github.com/PuerkitoBio/goquery v1.9.2 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.10 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20251022180443-0feb69152e9f // indirect
	github.com/danieljoos/wincred v1.2.3 // indirect
	github.com/davidmz/go-pageant v1.0.2 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.35.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/go-fed/httpsig v1.1.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/godbus/dbus/v5 v5.2.0 // indirect
	github.com/google/go-github/v74 v74.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.11 // indirect
	github.com/googleapis/gax-go/v2 v2.16.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-version v1.8.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	github.com/tinylib/msgp v1.6.1 // indirect
	gitlab.com/gitlab-org/api/client-go v1.9.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.38.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/image v0.20.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/genproto v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120174246-409b4a993575 // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
	salsa.debian.org/vasudev/gospake2 v0.0.0-20210510093858-d91629950ad1 // indirect
)
