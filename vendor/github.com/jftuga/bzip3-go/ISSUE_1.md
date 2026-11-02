# Frame decoder trusted `orig_size` from block headers (out-of-range slice)

**Labels:** bug, security, fixed

## Summary

The frame decoder (`Decompress`) used the per-block `orig_size` field read
from the input stream to size a slice operation without first bounding it
against `Bound(blockSize)`. A crafted stream with an oversized `orig_size`
triggered a slice-out-of-range panic.

This was found by the native fuzz target `FuzzDecompressFrame`:

```
panic: runtime error: slice bounds out of range [:117440558] with capacity 67923
```

Here `67923` is the bound-sized block buffer and `117440558` is an
attacker-controlled `orig_size` taken directly from a mutated header.

## Reproduction

```
go test -run=xxx -fuzz=FuzzDecompressFrame -fuzztime=60s .
```

(Reproduced within seconds before the fix; the seed corpus mutates a valid
frame's block header.)

## Root cause

In the `Decompress` header pre-scan, `orig_size` was only checked for
`< 0`. A value larger than the block buffer's capacity then flowed into a
`buf[:origSize]` slice.

## Fix (applied)

1. In the pre-scan, reject any block whose `orig_size` exceeds
   `Bound(blockSize)` with `ErrMalformedHeader`.
2. After decoding each block, reject it if the actual decoded size differs
   from the advertised `orig_size` (see `ISSUE_3.md`).

The same hardening was applied to the streaming `Reader.nextBlock`.

## Upstream note (iczelia/bzip3)

Upstream `bz3_decompress` (`src/libbz3.c`) validates only `orig_size < 0`
at the frame level (line ~973) and relies on `bz3_decode_block`'s internal
`orig_size > bz3_bound(block_size)` check (line ~716) plus the `buf_max`
check to stay memory-safe. That layering appears correct in C, but the
frame-level validation is looser than the block decoder's. A defensive
upper-bound check at the header would make the trust boundary explicit.
Worth verifying and possibly reporting upstream.
