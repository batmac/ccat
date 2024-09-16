//rustc main.rs --target wasm32-wasip1 -o rust.wasm

use std::io::{self, Read, Write};

fn main() {
    let mut buffer = Vec::new();
    let _ = io::stdin().read_to_end(&mut buffer);
    let _ = io::stdout().write_all(&buffer);
    ()
}
