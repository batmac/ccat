# `needed_header_size` uses `(model & bit) * 4`, over-requiring header bytes

**Labels:** upstream, correctness, wontfix (bug-compatible)

## Summary

Upstream's block-size validation for the optional LZP/RLE size fields
computes the required header length with a bit value rather than a boolean:

```c
// src/libbz3.c, bz3_decode_block and bz3_orig_size_sufficient_for_decode
size_t needed_header_size = 9 + ((model & 2) * 4) + ((model & 4) * 4);
```

`model & 2` is `2` when the LZP bit is set (not `1`), so the LZP term
contributes `8` bytes instead of `4`; likewise `model & 4` contributes `16`
instead of `4`. The check therefore requires more header bytes than the
format actually uses (each optional size field is a single 4-byte value).

## Impact

Minor. It only over-estimates the minimum buffer size for the header, so it
can reject a buffer that is a few bytes smaller than strictly necessary. It
does not cause incorrect decoding or memory unsafety, because the actual
field reads use the correct `9 + 4 * p` offsets.

## Port status

The Go port preserves this behavior deliberately so that both
implementations accept and reject exactly the same inputs
(`bzip3.go`, `DecodeBlock`):

```go
// Bug-compatible with upstream: (model&2)*4 and (model&4)*4 evaluate
// to 8 and 16, over-requiring the header size slightly.
needed := 9 + int(model&modelLZP)*4 + int(model&modelRLE)*4
```

## Suggested upstream fix

If upstream chooses to correct it (which would be a format-observable
behavior change for pathologically small buffers):

```c
size_t needed_header_size = 9 + (!!(model & 2)) * 4 + (!!(model & 4)) * 4;
```
