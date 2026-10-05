All four optimization tasks from PERFORMANCE_IMPROVEMENT.md are implemented, verified, and committed on `dev` (5 commits, one per task plus docs). The headline: **decode went from ~20% slower than C bzip3 to parity, and encode is now faster than C**, with compressed output still byte-for-byte identical in both directions.

Measured results (Apple M4, shakespeare.txt benchmark, `-benchtime 10x`):

| Change | Encode | Decode |
|---|---|---|
| Baseline | 17.7 MB/s | 16.8 MB/s |
| 1. Bigram unBWT (libsais bucket2+fastbits port) | — | 23.1 MB/s |
| 2. cm coder locals + folded counter updates | 17.7 | 23.6 |
| 3. Slicing-by-8 CRC | 18.2 | 24.6 |
| 4. sais `int` → `int32` (memory only) | 18.0 | 24.4 |

On the 21 MiB single-block wall-clock test: Go encode 0.34 s user vs C 0.44 s; decode 0.24 s vs C 0.23 s. The int32 retype halved suffix-array alloc_space (512 → 256 MB in the profile), bringing encoder scratch to ~8x block size, same as C.

Notes from the work:

- Task 1 required two details beyond a straight port: the biPSI array needs `n+1` entries (so `lfBuf` grew by one), and it's zeroed per call so corrupt-but-header-valid input can't follow stale positions out of bounds — I traced the bound invariants and also fuzzed `FuzzDecodeBlock`/`FuzzRoundTrip` 30 s each with no findings.
- Adding the `unbwt` pointer field initially cost encode ~5% by shifting the embedded 128 KiB `cm.Coder` offset inside `State`; moving the field after `cm` restored the original layout.
- Task 2's `update0`/`update1` helpers were folded into direct stores of already-loaded values (provably overflow-free), confirmed bit-exact by the byte-identity tests.
- Task 4 was a mechanical word-boundary `int`→`int32` retype of both generated sais files; `ComputeSA` now enforces the 2^30 cap the algorithm's sign-bit flagging requires.

Every step was gated on `go test ./...` (69 tests), the C-interop byte-identity tests, and CLI cross-validation with C bzip3 1.5.3 in both directions including `-j` parallel mode — all passing at the end. README, PORTING_NOTES.md, and PERFORMANCE_IMPROVEMENT.md (marked completed) are updated with only measured numbers, and I saved a project memory so a future session doesn't redo this pass. Remaining CPU is now dominated by the arithmetic coder (~44% combined), which is bit-serial and can't be restructured without breaking bit-exactness — that's the wall for future tuning.
