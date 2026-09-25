package mutators

import (
	"fmt"
	"io"

	"github.com/batmac/ccat/pkg/log"
	"github.com/jftuga/bzip3-go"
)

func init() {
	singleRegister("unbzip3", unbzip3, withDescription("decompress bzip3 data"),
		withCategory("decompress"),
	)
	singleRegister("bzip3", cbzip3, withDescription("compress to bzip3 data (X:16777216 is block size in bytes, 65KiB-511MiB)"),
		withCategory("compress"),
		withConfigBuilder(bzip3BlockSizeConfig),
	)
}

func bzip3BlockSizeConfig(args []string) (any, error) {
	c, err := stdConfigHumanSizeAsInt64WithDefault(bzip3.DefaultBlockSize)(args)
	if err != nil {
		return nil, err
	}
	blockSize := c.(int64)
	if blockSize < bzip3.MinBlockSize || blockSize > bzip3.MaxBlockSize {
		return nil, fmt.Errorf("bzip3 block size %d out of range [%d, %d]", blockSize, bzip3.MinBlockSize, bzip3.MaxBlockSize)
	}
	return int32(blockSize), nil // #nosec G115 -- bounds checked above
}

func unbzip3(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return io.Copy(out, bzip3.NewReader(in))
}

func cbzip3(out io.WriteCloser, in io.ReadCloser, conf any) (int64, error) {
	blockSize := conf.(int32)
	log.Debugf("bzip3 block size: %d\n", blockSize)

	e, err := bzip3.NewWriter(out, blockSize)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(e, in)
	if err != nil {
		return n, err
	}
	return n, e.Close()
}
