# Performance Improvement Task

**Status: completed 2026-07-02.** All four tasks below were implemented
(one commit each). Results: EncodeBlock 17.7 -> 17.9 MB/s, DecodeBlock
16.8 -> 24.4 MB/s (decode now at parity with C bzip3; encode faster than
C). Encoder scratch dropped from ~12x to ~8x block size. Remaining CPU
is dominated by the arithmetic coder (~44% combined) and sais (~34%).
Kept for reference; see README.md and PORTING_NOTES.md for current
numbers.

This document is a self-contained task input for a fresh Claude session. It
describes measured hotspots in `github.com/jftuga/bzip3-go` and the
optimization work to attempt, in priority order.

## Ground rules (do not skip)

1. **Correctness is non-negotiable.** After every change, the full test
   suite must pass and compressed output must remain byte-for-byte identical
   to C bzip3:
   ```
   go test ./...
   go test -v -run 'TestByteIdentical|TestCross' .   # needs C bzip3 installed
   ```
   Decode-side changes must not alter output; encode-side changes must not
   alter the compressed bytes (BWT output is unique, so any divergence is a
   regression).
2. **Measure, don't guess.** Re-profile before and after each change with the
   commands below. Keep changes that show a real improvement; revert ones
   that don't.
3. Do not add cgo or third-party runtime dependencies.

## Current baseline (Apple M4, 2026-07-02)

Throughput is already competitive with C bzip3 1.5.3 — there is **no**
order-of-magnitude gap. A 21 MiB single-block encode measured 0.37s (Go) vs
0.45s (C) user time; decode is within ~20% of C. So this is incremental
tuning, not a rescue.

Benchmark harness (single 5.4 MiB block, shakespeare.txt):

```
go test -run xxx -bench 'BenchmarkEncodeBlock|BenchmarkDecodeBlock' \
  -benchtime 5x -cpuprofile cpu.prof -memprofile mem.prof .
go tool pprof -top -nodecount=20 cpu.prof
go tool pprof -top -sample_index=alloc_space mem.prof
```

Baseline: EncodeBlock ~17.7 MB/s, DecodeBlock ~16.8 MB/s.

## Profiled CPU hotspots (combined encode+decode run)

| Rank | Function | flat % | Notes |
|------|----------|--------|-------|
| 1 | `internal/bwt.Decode` | 22.0 | Naive LF-mapping unBWT. **Top target.** |
| 2 | `internal/cm.(*Coder).Encode` | 18.4 | Arithmetic coder inner loop |
| 3 | `internal/cm.(*Coder).Decode` | 14.0 | Arithmetic coder inner loop |
| 4 | `internal/sais.*` (computeSA etc.) | ~20 cum | Already competitive with libsais |
| 5 | `internal/crc.Sum` | 3.1 | Byte-at-a-time table CRC |

## Optimization targets, in priority order

### 1. `internal/bwt.Decode` — the biggest single win (decode)

The current inverse BWT is a textbook LF-mapping walk (one `lf[]` lookup and
one output byte per position, walking backward). It is the largest single
CPU consumer on the decode path at ~22%.

libsais's `libsais_unbwt_*` (in the reference repo's `include/libsais.h`,
functions around `libsais_unbwt_decode_*` and
`libsais_unbwt_calculate_biPSI`) is dramatically faster because it decodes
**bigrams** (two bytes per step) using a `bucket2` table of 256×256 entries
plus a `fastbits` index, halving the number of dependent-load iterations.

Task: port the bigram/`bucket2`+`fastbits` inverse-BWT strategy into
`internal/bwt.Decode`, keeping the existing libsais output convention
(`U[0]=T[n-1]`, primary index = skipped-row position + 1). Verify against
`TestKnownVector` and the round-trip and C-interop tests. Watch the extra
memory: `bucket2` is 256×256 `uint32` (256 KiB), acceptable and allocatable
once per `State`.

### 2. `internal/cm` — arithmetic coder inner loop (encode + decode, ~32%)

The per-bit loop indexes `s.c1[c1][ctx]` and `s.c2[2*ctx+f][j]` repeatedly.
Things to try and measure individually:

- Hoist `&s.c1[c1]` and `&s.c2[2*ctx+f]` into local slice variables to cut
  repeated bounds checks and multi-dimensional index math.
- Confirm `update0`/`update1` stay inlined (they currently do).
- Consider replacing the four separate `s.c2[...][j]` / `[j+1]` reads with a
  local pair.

The coder must remain bit-exact; the round-trip and byte-identical tests are
the guard. Do not restructure the probability math.

### 3. `internal/crc.Sum` — slice-by-8 (minor, ~3%)

Replace the byte-at-a-time loop with a slice-by-8 table (8 KiB of tables) as
`hash/crc32` does internally. Note bzip3 uses init value 1 and no final XOR;
keep `Sum(crc, buf)` semantics identical. Low priority but easy and safe.

### 4. Suffix array width: `[]int` -> `int32` (memory, not speed)

`internal/sais` uses `[]int` (8 bytes/entry). Since block sizes are capped at
511 MiB (< 2^31), the suffix array fits in `int32`. Converting halves the
encoder's SA scratch (`State.saBuf`, currently the largest allocation at
~66% of alloc_space in the benchmark). The vendored package is generated from
`sais_gen.go`; the cleanest route is to regenerate/retype the `byte`-input
variant to emit `int32` indices, or maintain a small local `int32` fork.
This is a memory optimization; do not expect a throughput change, and
confirm there is none.

## Reporting

When done, update the "Performance" section of `README.md` and the
`bzip3-go-port` memory with the new measured numbers. Do not claim
improvements you did not measure.
