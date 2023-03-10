package mutators

import (
	"io"
	"log"

	"github.com/batmac/ccat/pkg/shaped"
)

func init() {
	singleRegister("maxbw", maxbw,
		withDescription("limit the bandwidth to the specified value (bytes per second)"),
		withConfigBuilder(stdConfigHumanSizeAsInt64),
	)
}

func maxbw(w io.WriteCloser, r io.ReadCloser, conf any) (int64, error) {
	bw := conf.(int64)
	if bw <= 0 {
		log.Fatalln("max bandwidth must be greater than 0")
		return 0, nil
	}
	stream := shaped.NewReader(r, int(bw))
	return io.Copy(w, stream)
}
