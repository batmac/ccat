# Decoded block size is not validated against advertised `orig_size`

**Labels:** upstream, correctness, hardened-in-port

## Summary

When decoding a frame, upstream copies exactly `orig_size` bytes out of the
block buffer after `bz3_decode_block` returns, without checking that the
block actually decoded to `orig_size` bytes:

```c
// src/libbz3.c, bz3_decompress
bz3_decode_block(state, compression_buf, compression_buf_size, size, orig_size);
if (bz3_last_error(state) != BZ3_OK) { ... return last_error; }
memcpy(out + *out_size, compression_buf, orig_size);   // trusts orig_size
```

For a raw/stored block (`bwt_idx == -1`), the decoder returns
`compressed_size - 8`, which is derived from the stream and is not
necessarily equal to the `orig_size` recorded in the frame's block header.
A crafted frame can therefore make the output contain stale bytes from the
block buffer beyond what was actually decoded.

## Impact

Not a memory-safety issue (the copy length is bounded by earlier checks),
but a correctness / minor information-exposure concern: bytes never present
in a legitimate stream can appear in the decompressed output.

## Port status (hardened)

The Go port rejects any block whose decoded size disagrees with the
advertised `orig_size`, in both the frame decoder and the streaming reader:

```go
dSize, err := state.DecodeBlock(buf, size, origSize)
if err != nil { return nil, err }
if dSize != origSize {
    return nil, ErrMalformedHeader
}
```

This is one of the two intentional divergences from C documented in the
README.

## Suggested upstream fix

Compare the return value of `bz3_decode_block` against `orig_size` and
treat a mismatch as `BZ3_ERR_MALFORMED_HEADER`.
