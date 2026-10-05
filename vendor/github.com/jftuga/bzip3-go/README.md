# bzip3-go

A pure-Go port of the [bzip3](https://github.com/iczelia/bzip3) compression
algorithm by Kamila Szewczyk. No cgo. Fully file-format compatible with the
C implementation, and verified to produce bit-for-bit identical compressed
output.

The codec pipeline is, per independent block:

1. RLE - per-byte-value run collapsing, applied only when it shrinks the block
2. LZP - hash-table match prediction (2^18 entries, minimum match 40), applied only when it shrinks the block
3. BWT - Burrows-Wheeler transform via a linear-time SA-IS suffix array
4. An order-0 context-mixing arithmetic coder

## Disclaimer

This port was vibe-coded with Claude Fable 5 as a fun experiment to see
how well an AI-assisted port of a C library and its CLI would turn out.
In my limited testing it produced byte-for-byte identical artifacts to
the C implementation, but it has not seen production use. Anyone who
depends on bzip3 for real work should use the reference C implementation
at <https://github.com/iczelia/bzip3>.

## Library

```go
import bzip3 "github.com/jftuga/bzip3-go"
```

Streaming (the bzip3 file format, same as the `bzip3` CLI tool):

```go
w, _ := bzip3.NewWriter(dst, bzip3.DefaultBlockSize)
io.Copy(w, src)
w.Close()

r := bzip3.NewReader(src)
io.Copy(dst, r)
```

One-shot frame API (equivalent to upstream `bz3_compress`/`bz3_decompress`):

```go
compressed, err := bzip3.Compress(data, bzip3.DefaultBlockSize)
original, err := bzip3.Decompress(compressed)
```

Low-level block API (equivalent to `bz3_encode_block`/`bz3_decode_block`),
for parallel or custom-container use:

```go
state, _ := bzip3.NewState(blockSize)
n, err := state.EncodeBlock(buf, dataSize) // in place; len(buf) >= bzip3.Bound(dataSize)
n, err := state.DecodeBlock(buf, compressedSize, origSize)
```

A `State` is not safe for concurrent use; create one per goroutine.

## CLI

```console
$ go install github.com/jftuga/bzip3-go/cmd/bzip3@latest

$ bzip3 -e file          # -> file.bz3
$ bzip3 -d file.bz3      # -> file
$ bzip3 -t file.bz3      # integrity test
$ bzip3 -e -b 32 -j 8 -c < input > output.bz3   # 32 MiB blocks, 8 workers
```

Flags mirror the C tool: `-e`/`-z` encode, `-d` decode, `-t` test, `-b`
block size in MiB (default 16), `-j` parallel workers, `-c` standard
streams, `-f` force overwrite, `-v` verbose, `-rm` remove input, `-V`
version. Note Go's flag parsing does not support combined short options
(use `-f -e -b 6`, not `-feb 6`).

## Compatibility and testing

- Decodes streams produced by C bzip3, and C bzip3 decodes streams produced
  by this port (verified in CI against the Debian `bzip3` package).
- Compressed output is byte-identical to C bzip3 1.5.3 for both the
  sequential and parallel paths. BWT output is unique for a given input and
  every other stage is deterministic, so this is expected to hold for all
  inputs; round-trip integrity is what the test suite guarantees.
- Upstream's test targets are ported: decoding `shakespeare.txt.bz3` (made
  by the C implementation) and the LICENSE round trip.
- Native Go fuzz targets port upstream's AFL harnesses (round trip, frame
  decompression, block decoding against adversarial input).

## Performance

On an Apple M4 (Go 1.26), throughput matches or beats C bzip3 1.5.3. On a
21 MiB single-block input, Go encode measured 0.34s user vs 0.44s for C,
and decode 0.24s vs 0.23s (parity). The package benchmarks on the 5.2 MiB
shakespeare.txt block measure ~17.9 MB/s encode and ~24.4 MB/s decode,
up from 17.7 and 16.8 before the inverse BWT was switched to libsais's
bigram (bucket2 + fastbits) strategy and the arithmetic-coder and CRC
inner loops were tuned. These are single-machine measurements, not a
rigorous benchmark. See `PERFORMANCE_IMPROVEMENT.md` for the profiled
hotspots this work was based on.

## Differences from the C implementation

- Malformed streams whose blocks decode to a size other than the advertised
  original size are rejected (`ErrMalformedHeader`); upstream copies stale
  buffer bytes in some of these cases.
- `bz3_recover` mode and the batch (`-B`) CLI mode are not implemented.

## License

LGPL-3.0-or-later, as this is a derivative work of libbz3 (Copyright
Kamila Szewczyk). The SA-IS suffix array construction in `internal/sais`
is adapted from [dsnet/compress](https://github.com/dsnet/compress)
(BSD 3-clause, Joe Tsai), itself a port of sais-lite by Yuta Mori (MIT);
see `internal/sais/LICENSE.md`.
