package openers

/*
import (
	"bytes"
	"ccat/log"
	"context"
	"fmt"
	"io"

	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/afs/scp"
)

var ScpOpenerName = "scp"
var ScpOpenerDescription = "get URL via SCP (afs)"

type ScpOpener struct {
	name, description string
}

func init() {
	register(&ScpOpener{
		name:        ScpOpenerName,
		description: ScpOpenerDescription,
	})
}

func (f ScpOpener) Name() string {
	return f.name
}
func (f ScpOpener) Description() string {
	return f.description
}
func (f ScpOpener) Open(s string, _ bool) (io.ReadCloser, error) {
	auth, err := scp.LocalhostKeyAuth("/User/bat/.ssh/id_ed25519")

	if err != nil {
		log.Fatal(err)
	}
	service := afs.New()
	ctx := context.Background()
	r, err := service.DownloadWithURL(ctx, s, auth, option.NewTimeout(2000))
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(r)
	//defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("data: %s\n", data)
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (f ScpOpener) Evaluate(s string) float32 {

	return 0.5
}
*/
