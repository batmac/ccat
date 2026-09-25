package mutators

import (
	"io"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy/xerial"
)

// Snappy/S2 formats that are not the framed stream handled by snap/s2.
// None of them is streamable: the whole input is read first.
func init() {
	// raw blocks: a varint of the decoded length, then the compressed data
	// (LevelDB/RocksDB blocks, Parquet pages, RPC payloads, ...).
	// The S2 decoder also reads Snappy blocks, so one decoder serves both.
	singleRegister("unsnapblock", unsnapblock, withDescription("decompress a raw (unframed) snappy or s2 block"),
		withCategory("decompress"),
		withAliases("uns2block"),
	)
	singleRegister("snapblock", csnapblock, withDescription("compress to a raw (unframed) snappy block"),
		withCategory("compress"),
	)
	singleRegister("s2block", cs2block, withDescription("compress to a raw (unframed) s2 block"),
		withCategory("compress"),
	)

	// xerial framing, from snappy-java (used by Kafka): not compatible with
	// the official snappy framing that unsnap reads.
	singleRegister("unxerial", unxerial, withDescription("decompress xerial-framed snappy data (snappy-java, Kafka), or a raw snappy block"),
		withCategory("decompress"),
	)
	singleRegister("xerial", cxerial, withDescription("compress to xerial-framed snappy data (snappy-java, Kafka)"),
		withCategory("compress"),
	)
}

// blockMutator applies a whole-buffer codec.
func blockMutator(out io.Writer, in io.Reader, codec func([]byte) ([]byte, error)) (int64, error) {
	src, err := io.ReadAll(in) // NOT streamable
	if err != nil {
		return 0, err
	}
	dst, err := codec(src)
	if err != nil {
		return 0, err
	}
	n, err := out.Write(dst)
	return int64(n), err
}

func unsnapblock(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return blockMutator(out, in, func(src []byte) ([]byte, error) { return s2.Decode(nil, src) })
}

func csnapblock(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return blockMutator(out, in, func(src []byte) ([]byte, error) { return s2.EncodeSnappy(nil, src), nil })
}

func cs2block(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return blockMutator(out, in, func(src []byte) ([]byte, error) { return s2.Encode(nil, src), nil })
}

func unxerial(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return blockMutator(out, in, xerial.Decode)
}

func cxerial(out io.WriteCloser, in io.ReadCloser, _ any) (int64, error) {
	return blockMutator(out, in, func(src []byte) ([]byte, error) { return xerial.Encode(nil, src), nil })
}
