# ccat

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
[![Go](https://github.com/batmac/ccat/actions/workflows/go.yml/badge.svg)](https://github.com/batmac/ccat/actions/workflows/go.yml)
![GitHub](https://img.shields.io/github/license/batmac/ccat)
[![Go Report Card](https://goreportcard.com/badge/github.com/batmac/ccat)](https://goreportcard.com/report/github.com/batmac/ccat)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbatmac%2Fccat.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbatmac%2Fccat?ref=badge_shield)
[![codecov](https://codecov.io/gh/batmac/ccat/branch/main/graph/badge.svg?token=PCD6DM6S75)](https://codecov.io/gh/batmac/ccat)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/batmac/ccat/main.svg)](https://results.pre-commit.ci/latest/github/batmac/ccat/main)

cat on steroids.
Leveraging great go modules to ease my CLI life.

## Install

### Homebrew

`brew install batmac/tap/ccatos`

### Manually

Get the latest release from <https://github.com/batmac/ccat/releases/latest>

### Build from source with Mage

Run the zero-install code or use your mage binary:

```shell
$ git clone https://github.com/batmac/ccat
$ go run magefiles/mage.go -l # or 'mage -l'
Targets:
  buildDefault
  buildDefaultAndTest*    buildDefault,test
  buildFull               tags: libcurl,crappy,plugins
  buildMinimal            tags: nohl,fileonly
  clean
  install                 put ccat to $GOPATH/bin/ccat
  installDeps             go mod download
  test                    all
  testCompressionGo
  testGo                  go test ./...
  updateREADME
  verifyDeps              go mod verify

* default target
$ go run magefiles/mage.go # or 'mage'
```

## Update

- You can update to the latest github release with `ccat --selfupdate`.
- You can check against your current installed version with `ccat --check`.

## Build tags

available build tags:

- `libcurl`: build with the libcurl opener.
- `libcurl_purego`: (macOS only) build with the libcurl opener without cgo: the system libcurl is loaded at runtime with [purego](https://github.com/ebitengine/purego). Ignored if `libcurl` is also set with cgo enabled.
- `plugins`: build with the yaegi plugins engine.
- `fileonly`: build with the local file opener only.
- `nohl`: build without the syntax-highlighter.
- `crappy`: build with some crappy (but useful) openers/mutators.
- `keystore`: build with the OS keyring support (Mac,Linux,Windows).

## Examples

```shell
$ ccat -m "x:kubectl get all" -i -w -t ready,Running --bg
NAME           READY   STATUS    RESTARTS   AGE
pod/busybox3   1/1     Running   1          17m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   21d
```

( "READY and "Running" have different background colors )

```shell
$ echo Hello: | ccat -m y2j,j2y,base64,hexdump
00000000  53 47 56 73 62 47 38 36  49 47 35 31 62 47 77 4b  |SGVsbG86IG51bGwK|
```

```shell
$ ccat -m zstd,base64 LICENSE
KLUv/WQwA+UUAGKyjCQQi1gA+Mu1uPH4trE1CdEhIptwP1oghhOOOP5rTRkCAKACAAdroEXGmA2Uz10inntYY97kIsI70zHae3CBb+DnAi4lAWNzxHPSRXalSceYEc4rc2FxPri8gX/QzYX9sT/32JzZHMbds2JzJmwIlmDMlYRjz/LGa8+BueftcBkb80Dn00G42FDA5RFeRujMBVvNmX8e5sjP1xLa4ftzH5ivJRriudg89+A789ibhIq4UDyyZ7G4GZni7BE2l7K96ZqlNKZ73IHNx5wltsQbeNlcc+bhPNOmS7jnE77PxeGI55jFszx3hWnuUSTcivBStqhQaFeY5xDeeOwNfGM+FiNjfC6gMx3jCoB0lwzgTGNcAoyKS/cWGWM2UL45c60B6ko+0c6WZuyo8Cl6249qdXyaqIRxSxnqfO0rw3XQV2mudYDt7dyRvdAhthvuNN3Hzc9QXbj4N6oSbktVj4U7cqWG33ycO0JVyW/+cW3JhdY0VXzqEkSWXpkcFb5Sc+VoKuOOrki4J0clqZErk0dHN3JFqbj6zQu/SrkjPUFsR9UE2b5ib2eoqYwr9VxdKXfDX6V822uiNRIAkACvUg44rlM/bm9nmm2NVNCNHtVMG7oY00eWvlK0IT8uPm4p/3HXRHKh00Ns26bZcW8vpqrkE0t/85UKoagE/XFfK/W4EtVsKQX5xFGpKh337W2MO0JVxVr2ZkuX14w7S+jhBgYFBheAAp4lSFMTyV3WqthWqJ9PlwgmAGYTiQo1Zi5cSUcnOBoEX3ev5YUMcVnvUrgx5HaCOLWhPcddUjOGZmX2yZFNAyB+ELmdEtsBPE17UsXU4fTV1jdNHpPl8KiAkC34tsZqWhxaa3RlIgn25jJkslDF+VBC8cW5mikplMUF0g==
```

## Docker image

multi-arch container images with tags `libcurl,crappy,plugins` are automatically built by [Github Actions](https://github.com/batmac/ccat/actions/workflows/docker-images.yml) and made available on [Dockerhub](https://hub.docker.com/r/batmac/ccat) and [Github Packages](https://github.com/batmac/ccat/pkgs/container/ccat). <br/>
for instance:

```shell
$ docker run --rm ghcr.io/batmac/ccat:latest -h
...
```

or

```shell
$ kubectl run -i --tty ccat --image=batmac/ccat:latest -- /bin/sh
...
```

## help

```
version v1.20.1-252-gecf0872b [libcurl,crappy,plugins,keystore,gcp,aws], commit ecf0872bc7b3df5dd51c5994464f06d48b6d30a0, built at 2026-09-25@21:20:31+0000 by Mage (go1.27.0 linux/amd64)
usage: ccat [options] [file ...]
  -t, --tokens string       comma-separated list of tokens
  -i, --ignore-case         tokens given with -t are case-insensitive
  -o, --only                don't display lines without at least one token
  -r, --raw                 don't treat tokens as regexps
  -n, --line-number         number the output lines, starting at 1.
  -L, --flock-in            exclusively flock each file before reading
  -l, --flock-out           exclusively flock stdout
  -w, --word                read word by word instead of line by line (only works with utf8)
  -b, --bg                  colorize the background instead of the font
  -H, --humanize            try to do what is needed to help (syntax-highlight, autodetect, etc.)
  -S, --style string        style to use (only used if -H, --fullhelp for the list)
  -F, --formatter string    formatter to use (only used if -H, --fullhelp for the list)
  -P, --lexer string        lexer to use (only used if -H, --fullhelp for the list)
  -m, --mutators string     mutators to use (comma-separated), --fullhelp for the list
  -V, --version             print version on stdout
      --license             print license on stdout
  -B, --buildinfo           print build info on stdout
  -h, --help                print usage on stderr
      --fullhelp            print full usage on stdout
      --selfupdate          Update to latest Github release
      --check               Check version with the latest Github release
      --forceupdate         Force overwriting to the latest Github release
  -d, --debug               debug what we are doing
  -M, --mem-usage           print memory usage on stderr at the end
  -k, --insecure            get files insecurely (globally)
  -C, --completion string   print shell completion script
  -T, --ui                  display with a minimal ui
      --pprof               enable cpu and mem profiling
      --setkey              interactively ask and store a secret in the OS keyring

---
ccat <files>...
 - highlighter (used with -H):
  - Lexers: 1S, 1S:Enterprise, ABAP, abap, abl, ABNF, abnf, aconf, actionscript, ActionScript, ActionScript 3, actionscript3, Ada, ada, ada2005, ada95, Agda, agda, ahk, AL, al, alloy, Alloy, ampl, AMPL, Angular2, ansible, antlr, ANTLR, apache, apacheconf, ApacheConf, apl, APL, AppleScript, applescript, aql, ArangoDB AQL, Arduino, arduino, arexx, armasm, ArmAsm, art, arturo, Arturo, as, as3, asm, ATL, atl, AutoHotkey, autohotkey, AutoIt, autoit, awk, Awk, b3d, Ballerina, ballerina, bash, Bash, Bash Session, bash-session, basic, bat, batch, Batchfile, Beef, beef, bf, bib, BibTeX, bibtex, bicep, Bicep, bind, BlitzBasic, blitzbasic, BNF, bnf, bplus, BQN, bqn, brainfuck, Brainfuck, bsdmake, c, C, c#, C#, c++, C++, C3, c3, caddy, caddy-d, caddyfile, Caddyfile, Caddyfile Directives, caddyfile-d, caddyfile-directives, Cap'n Proto, capnp, cassandra, Cassandra CQL, cassette, Ceylon, ceylon, cf3, cfengine3, CFEngine3, cfg, cfs, cfstatement, chai, ChaiScript, chaiscript, Chapel, chapel, Cheetah, cheetah, chpl, cl, clj, clojure, Clojure, cmake, CMake, cobol, COBOL, coffee, coffee-script, CoffeeScript, coffeescript, Common Lisp, common-lisp, console, containerfile, Coq, coq, Core, core, cpp, cql, cr, crystal, Crystal, csh, csharp, CSS, css, CSV, csv, cucumber, Cucumber, CUE, cue, Cython, cython, d, D, dart, Dart, Dax, dax, desktop, Desktop file, desktop_entry, devicetree, Devicetree, Diff, diff, django, Django/Jinja, dns, docker, Docker, dockerfile, dosbatch, dosini, DTD, dtd, dts, duby, Dylan, dylan, ebnf, EBNF, edn, elisp, elixir, Elixir, Elm, elm, emacs, emacs-lisp, EmacsLisp, erb, ERB, Erlang, erlang, ex, exs, f90, Factor, factor, Fennel, fennel, fish, Fish, fishshell, fnl, Forth, forth, fortran, Fortran, FortranFixed, fortranfixed, fsharp, FSharp, GAS, gas, gawk, gd, gd3, GDScript, gdscript, GDScript3, gdscript3, gemfile-lock, Gemfile.lock, gemfilelock, gemini, Gemtext, gemtext, Genshi, genshi, Genshi HTML, Genshi Text, genshitext, Gettext, Gherkin, Gherkin, gherkin, Gleam, gleam, GLSL, glsl, gmi, gmni, Gnuplot, gnuplot, go, Go, Go HTML Template, Go Template, Go Text Template, go-html-template, go-template, go-text-template, golang, gql, graphql, GraphQL, graphqls, groff, Groff, Groovy, groovy, gsed, Handlebars, handlebars, hare, Hare, Haskell, haskell, Haxe, haxe, hbs, HCL, hcl, hcl, Hexdump, hexdump, hlb, HLB, HLSL, hlsl, HolyC, holyc, hs, html, HTML, html+erb, html+genshi, html+kid, html+ruby, HTTP, http, hx, hxsl, Hy, hylang, idr, idris, Idris, Igor, igor, igorpro, ini, INI, Io, io, iscdhcpd, ISCdhcpd, J, j, janet, Janet, java, Java, java-properties, javascript, JavaScript, jinja, jl, js, json, JSON, JSONata, jsonata, jsonl, Jsonnet, jsonnet, jsx, Julia, julia, jungle, Jungle, kak, Kakoune, kakoune, kakrc, kakscript, KDL, kdl, kid, kotlin, Kotlin, ksh, lateralus, Lateralus, latex, lean, lean4, Lean4, lighttpd, Lighttpd configuration file, lighty, lilypond, LilyPond, lisp, LLVM, llvm, lox, ltl, lua, Lua, luau, Luau, m2, make, Makefile, makefile, Mako, mako, man, mariadb, markdown, Markless, mason, Mason, materialize, Materialize SQL dialect, mathematica, Mathematica, matlab, Matlab, mawk, mbt, mcf, MCFunction, mcfunction, md, Meson, meson, meson.build, mess, metal, Metal, mf, microcad, minizinc, MiniZinc, mkd, MLIR, mlir, mma, Modelica, modelica, Modula-2, modula2, mojo, Mojo, MonkeyC, monkeyc, moon, moonbit, MoonBit, moonscript, MoonScript, morrowind, MorrowindScript, mwscript, Myghty, myghty, mysql, MySQL, MZN, mzn, mzsql, NASM, nasm, natural, Natural, nawk, nb, NDISASM, ndisasm, newspeak, Newspeak, ng2, nginx, Nginx configuration file, Nim, nim, nimrod, Nix, nix, nixos, no-highlight, nroff, nsh, nsi, NSIS, nsis, Nu, nu, obj-c, objc, Objective-C, objective-c, objectivec, objectpascal, ObjectPascal, OCaml, ocaml, Octave, octave, odin, Odin, ones, onesenterprise, OnesEnterprise, openedge, OpenEdge ABL, openedgeabl, OpenSCAD, openscad, org, Org Mode, orgmode, PacmanConf, pacmanconf, perl, Perl, perl6, PHP, php, php3, php4, php5, phtml, PHTML, Pig, pig, PkgConfig, pkgconfig, pl, PL/pgSQL, pl6, plain, plaintext, plc, plpgsql, Plutus Core, plutus-core, po, Pony, pony, posh, postgres, postgresql, PostgreSQL SQL dialect, postscr, postscript, PostScript, pot, pov, POVRay, PowerQuery, powerquery, PowerShell, powershell, pq, progress, prolog, Prolog, Promela, promela, promql, PromQL, properties, proto, protobuf, Protocol Buffer, Protocol Buffer Text Format, prql, PRQL, ps1, psd1, psl, PSL, psm1, puppet, Puppet, pwsh, py, py2, py3, pyrex, python, Python, Python 2, python2, python3, pyx, QBasic, qbasic, qbs, qml, QML, R, r, racket, Racket, Ragel, ragel, raku, Raku, rb, react, react, reason, ReasonML, reasonml, reg, registry, rego, Rego, rest, reStructuredText, restructuredtext, rexx, Rexx, rgbasm, RGBDS Assembly, rhtml, Ring, ring, rkt, RPG IV, RPGLE, RPMSpec, rs, rst, Ruby, ruby, rust, Rust, s, sage, salt, SAS, sas, Sass, sass, Scala, scala, scdoc, scdoc, scheme, Scheme, Scilab, scilab, scm, SCSS, scss, Sed, sed, sh, shell, shell-session, Sieve, sieve, sls, Smali, smali, smalltalk, Smalltalk, Smarty, smarty, sml, snbt, SNBT, snobol, Snobol, sol, solidity, Solidity, SourcePawn, sp, Spade, spade, sparql, SPARQL, spec, spitfire, splus, sql, SQL, SQLRPGLE, squeak, squid, squid.conf, squidconf, SquidConf, ssed, st, Standard ML, starlark, stas, Stylus, stylus, sv, svelte, Svelte, swift, Swift, SYSTEMD, systemd, systemverilog, systemverilog, t-sql, tablegen, TableGen, tal, Tal, tape, TASM, tasm, Tcl, tcl, tcsh, Tcsh, templ, Templ, Termcap, termcap, Terminfo, terminfo, Terraform, terraform, TeX, tex, text, tf, thrift, Thrift, TOML, toml, TradingView, tradingview, Transact-SQL, ts, tsql, tsx, Turing, turing, Turtle, turtle, tv, Twig, twig, txtpb, TypeScript, typescript, typoscript, TypoScript, typoscriptcssdata, TypoScriptCssData, typoscripthtmldata, TypoScriptHtmlData, typst, Typst, ucode, udiff, uxntal, v, V, v, V shell, Vala, vala, vapi, VB.net, vb.net, vbnet, verilog, verilog, vhdl, VHDL, VHS, vhs, vim, VimL, vlang, vsh, vshell, vtt, vue, vue, vuejs, wast, wat, WDTE, WebAssembly Text Format, WebGPU Shading Language, WebVTT, wgsl, Whiley, whiley, winbatch, xml, XML, xml+genshi, xml+kid, Xorg, xorg.conf, YAML, yaml, yaml+jinja, YAML+Jinja, YANG, yang, z80, Z80 Assembly, Zed, zed, Zig, zig, zone, zsh, µcad, 🔥
  - Styles: abap, algol, algol_nu, arduino, ashen, aura-theme-dark, aura-theme-dark-soft, autumn, average, base16-snazzy, borland, bw, catppuccin-frappe, catppuccin-latte, catppuccin-macchiato, catppuccin-mocha, colorful, darcula, doom-one, doom-one2, dracula, emacs, evergarden, friendly, fruity, github, github-dark, gruvbox, gruvbox-light, hr_high_contrast, hrdark, igor, kanagawa-dragon, kanagawa-lotus, kanagawa-wave, lovelace, manni, modus-operandi, modus-vivendi, monokai, monokailight, murphy, native, nord, nordic, onedark, onesenterprise, paraiso-dark, paraiso-light, pastie, perldoc, pygments, rainbow_dash, rose-pine, rose-pine-dawn, rose-pine-moon, rpgle, rrt, solarized-dark, solarized-dark256, solarized-light, swapoff, tango, tokyonight-day, tokyonight-moon, tokyonight-night, tokyonight-storm, trac, vim, vs, vulcan, witchhazel, xcode, xcode-dark
  - Formatters: html, json, noop, svg, terminal, terminal16, terminal16m, terminal256, terminal8, tokens
 - openers:
    cb: get content from the system clipboard via cb://
    crng: get data from crypto/rand (accept a size limit as parameter)
    echo: echo the string given
    file: open local files
    gcs: get a GCP Cloud Storage object via gs://
    gemini: get URL via Gemini
    http: get URL via HTTP(S)
    curl: get URL via libcurl bindings
           libcurl/8.5.0 OpenSSL/3.0.13 zlib/1.3 brotli/1.1.0 zstd/1.5.5 libidn2/2.3.7 libpsl/0.21.2 (+libidn2/2.3.7) libssh/0.10.6/openssl/zlib nghttp2/1.59.0 librtmp/2.3 OpenLDAP/2.6.10
           protocols: dict,file,ftp,ftps,gopher,gophers,http,https,imap,imaps,ldap,ldaps,mqtt,pop3,pop3s,rtmp,rtmpe,rtmps,rtmpt,rtmpte,rtmpts,rtsp,scp,sftp,smb,smbs,smtp,smtps,telnet,tftp
    mc: get a Minio-compatible object via mc:// (use ~/.mc/config.json or env for credentials)
    tcp: get data from listening on tcp://[HOST]:<PORT>
    prng: generate endless pcg rand (don't use for crypto) (accept a seed as parameter)
    s3: get an AWS s3 object via s3://
    ShellScp: get scp:// via local scp
    sse: stream Server-Sent Events via sse://
    wormhole: get text, file or zipped dir via a wormhole code (wh://<code> or wormhole://<code>)
 - mutators:
        cb: put a copy in the clipboard
        discard: discard X:0 bytes (0 = all)
        dummy: a simple fifo
        help: display mutators help
        hexdump: dump in hex as xxd
        indent: indent the text (with X:4 chars)
        j: JSON Re-indent (X:2 space-based)
        limit: a simple limiting fifo ( with X max size in bytes, for instance 'limit:1k')
        maxbw: limit the bandwidth to the specified value (bytes per second)
        mimetype: detect mimetype
        pv: copy in to out, printing the total and the bandwidth (like pv) each X:1000 milliseconds on stderr
        sponge: soak all input before outputting it.
        wc: count bytes (b, default), runes (r), words (w) or lines (l)
        wrap: word-wrap the text (to X:80 chars maximum)
        wrapU: unconditionally wrap the text (to X:80 chars maximum)
        x: execute command (e.g. 'x:head -n 10')
    checksum:
        crc32: compute the crc32 checksum
        md5: compute the md5 checksum
        sha1: compute the sha1 checksum
        sha224: compute the sha224 checksum
        sha256: compute the sha256 checksum
        sha3-224: compute the sha3-224 checksum
        sha3-256: compute the sha3-256 checksum
        sha3-384: compute the sha3-384 checksum
        sha3-512: compute the sha3-512 checksum
        sha384: compute the sha384 checksum
        sha512: compute the sha512 checksum
        xxh32: compute the xxhash32 checksum
        xxh3: compute the xxh3 checksum
        xxh64: compute the xxhash64 checksum
    compress:
        brotli: compress to brotli data (X:6 is compression level, 0-11)
        bzip2: compress to bzip2 data (X:9 is compression level, 0-9)
        bzip3: compress to bzip3 data (X:16777216 is block size in bytes, 65KiB-511MiB)
        deflate: compress to raw deflate data (RFC 1951, no header) (X:6 is compression level, 0-9)
        gzip: compress to gzip data (X:6 is compression level, 0-9)
        lz4: compress to lz4 data (X:0 is compression level, 0-9)
        lzma2: compress to lzma2 data
        lzma: compress to lzma data
        minlz: compress to minlz data
        pbzip3: parallel compress to bzip3 data (X:0 is concurrency, 0 is auto, then X:16777216 is block size in bytes)
        pgzip: compress with pgzip  (X:6 is compression level, 0-9, blockSize, blocks)
        s2: compress to s2 data
        s2block: compress to a raw (unframed) s2 block
        snap: compress to snappy data
        snapblock: compress to a raw (unframed) snappy block
        xerial: compress to xerial-framed snappy data (snappy-java, Kafka)
        xz: compress to xz data
        zip: compress to zip data
        zlib: compress to zlib data (X:6 is compression level, 0-9)
        zstd: compress to zstd data (X:4 is compression level, 1-22)
    convert:
        base64: encode to base64
        feed2y: rss/atom/json feed -> YAML
        hex: dump in lowercase hex
        html2md: html -> markdown
        it2dl: download via iTerm2 escape code, does not work in other terminals. Must be the last mutator of the pipeline.
        j2y: JSON -> YAML
        j5j: JSON5 -> JSON
        jcs: JSON -> JSON Canonicalization (RFC 8785)
        ndjsonindent: pretty-print NDJSON/JSON Lines
        plist2Y: display an Apple plist as yaml
        qp: encode quoted-printable data
        unbase64: decode base64
        unhex: decode hex, ignore all non-hex chars
        unqp: decode quoted-printable data
        y2j: YAML -> JSON
    decompress:
        punbzip3: parallel decompress bzip3 data (X:0 is concurrency, 0 is auto)
        unbrotli: decompress brotli data
        unbzip2: decompress bzip2 data in parallel (X:0 is concurrency, 0 is one worker per CPU)
        unbzip3: decompress bzip3 data
        undeflate: decompress raw deflate data (RFC 1951, no header)
        ungzip: decompress gzip data
        unlz4: decompress lz4 data
        unlzfse: decompress lzfse data
        unlzma2: decompress lzma2 data
        unlzma: decompress lzma data
        unminlz: decompress minlz data
        uns2: decompress s2 data
        unsnap: decompress snappy data
        unsnapblock: decompress a raw (unframed) snappy or s2 block
        unxerial: decompress xerial-framed snappy data (snappy-java, Kafka), or a raw snappy block
        unxz: decompress xz data
        unzip: decompress the first file in a zip archive
        unzlib: decompress zlib data
        unzstd: decompress zstd data
    decrypt:
        easyopen: decrypt with Nacl EasyOpen, get the key from env (KEY)
    encrypt:
        easyseal: encrypt with Nacl EasySeal, key used is printed on stderr
    external APIs:
        chatgpt: ask OpenAI ChatGPT, X:<unlimited> max replied tokens, the optional second arg is the model (Requires a valid key in $OPENAI_API_KEY, optional custom endpoint in $OPENAI_BASE_URL.)
        claude: ask Anthropic Claude, X:<unlimited> max replied tokens, optional second arg is the model, optional third arg is the preprompt (needs a valid key in $ANTHROPIC_API_KEY)
        googleai: googleai, X:gemini-1.5-flash is the model (Requires a valid key in $GOOGLE_API_KEY)
        huggingface: ask HuggingFace for simple tasks, optional args are model, max tokens, temperature (needs a valid key in $HUGGING_FACE_HUB_TOKEN, set HUGGING_FACE_ENDPOINT to use an Inference API endpoint)
        mistralai: ask MistralAI, X:<unlimited> max replied tokens, the optional second arg is the model (Requires a valid key in $MISTRAL_API_KEY)
        translate: translate to X:en or $TARGET_LANGUAGE with google translate (needs a valid key in $GOOGLE_API_KEY)
        wa: query wolfram alpha Short Answers API (APPID in $WA_APPID)
        wallm: query wolfram alpha LLM API (APPID in $WA_APPID)
        wasimple: query wolfram alpha Simple API (output is an image, APPID in $WA_APPID)
        waspoken: query wolfram alpha Spoken API (APPID in $WA_APPID)
    filter:
        filterUTF8: remove non-utf8
        jsonpath: a jsonpath expression to apply (on $, with all ',' replaced by '|', all ':' replaced by '£')
        removeANSI: remove ANSI codes
        shuffle: group the X:8-byte elements by byte position, so a following compressor sees runs instead of interleaved bytes
        unshuffle: reverse shuffle (X:8 must match the shuffle element size)
    plugin:
        wasm: a wasi (wasm) module to apply (path as first argument)
        yaegi: a yaegi script to apply (path as first argument, symbol as second argument)

  ('X:Y' means X is an argument with default value Y)

  mutator aliases:
    b64: base64
    cgpt: chatgpt
    d: discard
    dum, dumm: dummy
    gai: googleai
    hd, xxd: hexdump
    h2m, h2md: html2md
    hf: huggingface
    head, l: limit
    mime: mimetype
    mistral: mistralai
    unpbzip3: punbzip3
    ub64, unb64: unbase64
    punbzip2: unbzip2
    inflate: undeflate
    unpgzip: ungzip
    uns2block: unsnapblock
```
