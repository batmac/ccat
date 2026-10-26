# Design Notes

Implementation notes for the pure-Go bzip3 port, ported from the C
reference at https://github.com/iczelia/bzip3 (`src/libbz3.c`).

## What was built

- **Library (root package):** the in-place block codec
  (`State.EncodeBlock` / `State.DecodeBlock`, mirroring `bz3_encode_block` /
  `bz3_decode_block`), a one-shot frame API (`Compress` / `Decompress`,
  equivalent to upstream `bz3_compress` / `bz3_decompress`), and a streaming
  `Reader` / `io.WriteCloser` pair for the CLI file format. Error values map
  1:1 to the upstream `BZ3_ERR_*` codes.
- **Stages (`internal/`):**
  - `crc` — bzip3's CRC32 variant (Castagnoli polynomial, initial value 1,
    no final XOR).
  - `rle` — the `mrlec` / `mrled` run-length stage.
  - `lzp` — the LZP predictor (2^18-entry hash table, minimum match 40).
  - `cm` — the order-0 context-mixing arithmetic coder.
  - `bwt` — BWT / unBWT using libsais output conventions, over an SA-IS
    suffix array adapted from dsnet/compress (Yuta Mori's sais-lite).
  - `sais` — the vendored SA-IS suffix array construction.
- **CLI (`cmd/bzip3`):** upstream-compatible flags (`-e`/`-d`/`-t`, `-b`,
  `-j` with goroutine workers, `-c`, `-f`, `-v`, `-rm`, `-V`).
- **Tests:** per-stage round trips, upstream's shakespeare and LICENSE
  targets, cross-validation and byte-identity against the C binary (skipped
  when it is not installed), error paths, and Go-fuzz ports of the three AFL
  harnesses.

Compressed output is bit-for-bit identical to C bzip3 1.5.3, verified in
both directions and for both the sequential and parallel paths.

## Things worth knowing

- **Fuzzing found a real bug during development.** The frame decoder trusted
  `orig_size` from block headers; a crafted stream caused an out-of-range
  slice. Fixed with bounds checks. Upstream `bz3_decompress` has the same
  pattern and will read stale or out-of-bounds buffer bytes there. See
  `ISSUE_1.md`.
- **Two deliberate divergences from C** (documented in the README):
  malformed blocks whose decoded size disagrees with the advertised size are
  rejected instead of emitting stale bytes, and `-B` / recover mode are not
  implemented.
- **A bug-compatible quirk was preserved on purpose.** Upstream's
  `needed_header_size` check computes `(model & 2) * 4` = 8 (not 4); the port
  keeps this so both implementations accept and reject exactly the same
  inputs. See `ISSUE_2.md`.
- **BWT convention.** libsais output is `U[0] = T[n-1]`, the suffix-array row
  equal to 0 is skipped, and the primary index is that row's position + 1.
  The unBWT is a hand-derived LF-mapping over the sentinel formulation.
- **Suffix array width.** The vendored SA-IS uses `[]int` (8 bytes/entry on
  64-bit) rather than C's `int32_t`, so the encoder needs about
  `12 * blockSize` bytes of scratch versus C's ~8. This is a memory
  difference, not a speed one.

## Performance summary

Contrary to an earlier (incorrect) estimate, there is no order-of-magnitude
gap with C. On an Apple M4, a 21 MiB single-block encode measured 0.37s
(Go) vs 0.45s (C) user time; decode is within roughly 20%. Profiled
hotspots and optimization ideas are in `PERFORMANCE_IMPROVEMENT.md`.
