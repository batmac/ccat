// zig cc --target=wasm32-wasi main.zig -o zig.wasm
const std = @import("std");

pub fn main() !void {
    var stdin = std.io.getStdIn().reader();
    var stdout = std.io.getStdOut().writer();

    const bufferSize: usize = 4096;
    var buffer: [bufferSize]u8 = undefined;

    while (true) {
        const bytesRead = try stdin.read(buffer[0..]);
        if (bytesRead == 0) break;

        _ = try stdout.write(buffer[0..bytesRead]);
    }
}
