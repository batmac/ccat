# Porting Notes

Notes from the initial port of [bzip3](https://github.com/iczelia/bzip3)
(C, libbz3.c) to pure Go, completed 2026-07-01. Compressed output is
bit-for-bit identical to C bzip3 1.5.3, verified in both directions with
the reference binary, including parallel mode.

## What was built

- **Library** (root package): in-place block codec (`State.EncodeBlock` /
  `State.DecodeBlock` mirroring `bz3_encode_block` / `bz3_decode_block`),
  one-shot frame API (`Compress` / `Decompress` mirroring `bz3_compress` /
  `bz3_decompress`), and a streaming `Reader` / `io.WriteCloser` pair for
  the CLI file format. Errors map 1:1 to the `BZ3_ERR_*` codes.
- **Stages** in `internal/`:
  - `crc`: bzip3's CRC32 variant (Castagnoli polynomial, init value 1, no
    final XOR).
  - `rle`: the mrlec/mrled run-length stage.
  - `lzp`: the LZP predictor stage (2^18-entry hash table, minimum match 40).
  - `cm`: the order-0 context-mixing arithmetic coder.
  - `bwt`: BWT/unBWT with libsais conventions, over the SA-IS suffix array
    port in `internal/sais` (adapted from dsnet/compress, itself a port of
    Yuta Mori's sais-lite). The unBWT is an LF-mapping backward walk derived
    from the sentinel formulation.
- **CLI** `cmd/bzip3` with upstream-compatible flags (`-e`/`-d`/`-t`, `-b`,
  `-j` with goroutine worker batches, `-c`, `-f`, `-v`, `-rm`).
- **Tests**: 68 passing. Per-stage round-trips, upstream's shakespeare and
  LICENSE test targets, cross-validation and byte-identity tests against
  the C binary (skipped when not installed), error-path tests, and Go-fuzz
  ports of upstream's three AFL harnesses. CI covers Linux/macOS/Windows
  plus a cross-validation job against Debian's bzip3 package.

## Things worth knowing

- **Fuzzing found a real bug during development**: the frame decoder
  trusted `orig_size` from block headers; a crafted stream caused an
  out-of-range slice. Fixed with bounds checks. Upstream `bz3_decompress`
  has the same pattern on its raw-block path and can over-read its block
  buffer. See ISSUE_1.md.
- **Deliberate divergences** (also in README): malformed blocks whose
  decoded size disagrees with the advertised original size are rejected
  instead of emitting stale buffer bytes; `-B` batch mode and recover mode
  are not implemented (see ISSUE_3.md).
- **Preserved a bug-compatible quirk**: upstream's `needed_header_size`
  check computes `(model & 2) * 4` = 8 rather than the intended 4, so it
  over-requires the buffer size. Kept as-is so both implementations
  accept/reject identically. See ISSUE_2.md.
- **BWT conventions** (the key porting subtlety): libsais_bwt places the
  last input byte at output position 0, skips the suffix-array row whose
  value is 0, and returns that row's position plus one as the primary
  index. Equivalently, the output is the BWT of the sentinel-terminated
  string with the sentinel removed and the index marking where it was.
  Any correct suffix array yields this exact output, which is why
  byte-identity with libsais holds without porting libsais itself.
- **Arithmetic-coder flush bytes are not covered by any check**: flipping
  the final byte of a block often does not change the decoded output (true
  for the C implementation as well). Corruption tests must target the
  middle of the payload.
- **Performance** (Apple M4, Go 1.26.4, single thread, after the
  2026-07-02 optimization pass from PERFORMANCE_IMPROVEMENT.md): on a
  21 MiB single-block input, encode 0.34 s user vs 0.44 s for C (Go
  faster); decode 0.24 s vs 0.23 s (parity). The wins came from porting
  libsais's bigram unBWT (decode +36%), hoisting the arithmetic-coder
  counter rows into locals, and slicing-by-8 CRC. The sais port was also
  retyped to `int32`, so encoder scratch is ~8x block size, matching C.
