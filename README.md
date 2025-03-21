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
version v1.17.2-55-g8a8e91e7 [libcurl,crappy,plugins,keystore], commit 8a8e91e7fc2f9d88395cf52bb920ffa2df94afb4, built at 2025-03-20@12:41:09+0100 by Mage (go1.24.1 darwin/arm64)
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
  - Lexers: 1S, 1S:Enterprise, ABAP, abap, abl, ABNF, abnf, aconf, actionscript, ActionScript, ActionScript 3, actionscript3, Ada, ada, ada2005, ada95, Agda, agda, ahk, AL, al, alloy, Alloy, Angular2, antlr, ANTLR, apache, ApacheConf, apacheconf, APL, apl, AppleScript, applescript, aql, ArangoDB AQL, arduino, Arduino, arexx, ArmAsm, armasm, as, as3, asm, ATL, atl, autohotkey, AutoHotkey, autoit, AutoIt, Awk, awk, b3d, ballerina, Ballerina, Bash, bash, Bash Session, bash-session, basic, bat, batch, Batchfile, beef, Beef, bf, bib, bibtex, BibTeX, bicep, Bicep, bind, blitzbasic, BlitzBasic, BNF, bnf, bplus, bqn, BQN, brainfuck, Brainfuck, bsdmake, c, C, C#, c#, c++, C++, caddy, caddy-d, caddyfile, Caddyfile, Caddyfile Directives, caddyfile-d, caddyfile-directives, Cap'n Proto, capnp, cassandra, Cassandra CQL, cassette, Ceylon, ceylon, cf3, cfengine3, CFEngine3, cfg, cfs, cfstatement, chai, ChaiScript, chaiscript, chapel, Chapel, Cheetah, cheetah, chpl, cl, clj, Clojure, clojure, cmake, CMake, COBOL, cobol, coffee, coffee-script, coffeescript, CoffeeScript, Common Lisp, common-lisp, console, Coq, coq, cpp, cql, cr, crystal, Crystal, csh, csharp, CSS, css, csv, CSV, Cucumber, cucumber, cue, CUE, cython, Cython, d, D, dart, Dart, dax, Dax, desktop, Desktop file, desktop_entry, Diff, diff, django, Django/Jinja, dns, docker, Docker, dockerfile, dosbatch, dosini, dtd, DTD, duby, Dylan, dylan, EBNF, ebnf, edn, elisp, elixir, Elixir, elm, Elm, emacs, emacs-lisp, EmacsLisp, erlang, Erlang, ex, exs, f90, Factor, factor, Fennel, fennel, Fish, fish, fishshell, fnl, Forth, forth, fortran, Fortran, fortranfixed, FortranFixed, fsharp, FSharp, gas, GAS, gawk, gd, gd3, GDScript, gdscript, GDScript3, gdscript3, genshi, Genshi, Genshi HTML, Genshi Text, genshitext, gherkin, Gherkin, Gherkin, Gleam, gleam, glsl, GLSL, gnuplot, Gnuplot, go, Go, Go HTML Template, Go Template, Go Text Template, go-html-template, go-template, go-text-template, golang, gql, GraphQL, graphql, graphqls, Groff, groff, groovy, Groovy, gsed, Handlebars, handlebars, hare, Hare, Haskell, haskell, haxe, Haxe, hbs, hcl, HCL, hexdump, Hexdump, HLB, hlb, hlsl, HLSL, HolyC, holyc, hs, HTML, html, html+genshi, html+kid, http, HTTP, hx, hxsl, Hy, hylang, idr, idris, Idris, igor, Igor, igorpro, INI, ini, Io, io, iscdhcpd, ISCdhcpd, j, J, java, Java, java-properties, JavaScript, javascript, jinja, jl, js, JSON, json, JSONata, jsonata, Jsonnet, jsonnet, jsx, Julia, julia, jungle, Jungle, kid, kotlin, Kotlin, ksh, latex, lighttpd, Lighttpd configuration file, lighty, lisp, llvm, LLVM, lua, Lua, m2, make, Makefile, makefile, mako, Mako, man, mariadb, markdown, Mason, mason, materialize, Materialize SQL dialect, mathematica, Mathematica, matlab, Matlab, mawk, mcf, MCFunction, mcfunction, md, Meson, meson, meson.build, metal, Metal, mf, minizinc, MiniZinc, mkd, MLIR, mlir, mma, Modula-2, modula2, monkeyc, MonkeyC, morrowind, MorrowindScript, mwscript, Myghty, myghty, mysql, MySQL, mzn, MZN, mzsql, NASM, nasm, Natural, natural, nawk, nb, ndisasm, NDISASM, newspeak, Newspeak, ng2, nginx, Nginx configuration file, Nim, nim, nimrod, Nix, nix, nixos, no-highlight, nroff, nsh, nsi, NSIS, nsis, obj-c, objc, Objective-C, objective-c, objectivec, objectpascal, ObjectPascal, ocaml, OCaml, Octave, octave, Odin, odin, ones, OnesEnterprise, onesenterprise, openedge, OpenEdge ABL, openedgeabl, openscad, OpenSCAD, org, Org Mode, orgmode, PacmanConf, pacmanconf, Perl, perl, perl6, php, PHP, php3, php4, php5, phtml, PHTML, pig, Pig, pkgconfig, PkgConfig, pl, PL/pgSQL, pl6, plain, plaintext, plc, plpgsql, Plutus Core, plutus-core, pony, Pony, posh, postgres, postgresql, PostgreSQL SQL dialect, postscr, postscript, PostScript, pov, POVRay, powerquery, PowerQuery, powershell, PowerShell, pq, progress, prolog, Prolog, promela, Promela, PromQL, promql, properties, proto, protobuf, Protocol Buffer, PRQL, prql, ps1, psd1, PSL, psl, psm1, puppet, Puppet, pwsh, py, py2, py3, pyrex, python, Python, Python 2, python2, python3, pyx, QBasic, qbasic, qbs, qml, QML, r, R, racket, Racket, Ragel, ragel, Raku, raku, rb, react, react, reason, reasonml, ReasonML, reg, registry, Rego, rego, rest, restructuredtext, reStructuredText, rexx, Rexx, rkt, RPMSpec, rs, rst, ruby, Ruby, Rust, rust, s, sage, sas, SAS, sass, Sass, scala, Scala, scheme, Scheme, scilab, Scilab, scm, SCSS, scss, Sed, sed, sh, shell, shell-session, sieve, Sieve, smali, Smali, Smalltalk, smalltalk, smarty, Smarty, sml, snbt, SNBT, snobol, Snobol, sol, solidity, Solidity, SourcePawn, sp, sparql, SPARQL, spec, spitfire, splus, sql, SQL, squeak, squid, squid.conf, SquidConf, squidconf, ssed, st, Standard ML, stas, Stylus, stylus, sv, svelte, Svelte, swift, Swift, SYSTEMD, systemd, systemverilog, systemverilog, t-sql, tablegen, TableGen, tal, Tal, tape, tasm, TASM, Tcl, tcl, Tcsh, tcsh, Termcap, termcap, terminfo, Terminfo, Terraform, terraform, TeX, tex, text, tf, thrift, Thrift, TOML, toml, tradingview, TradingView, Transact-SQL, ts, tsql, tsx, Turing, turing, turtle, Turtle, tv, twig, Twig, typescript, TypeScript, TypoScript, typoscript, TypoScriptCssData, typoscriptcssdata, TypoScriptHtmlData, typoscripthtmldata, Typst, typst, ucode, udiff, uxntal, v, V, v, V shell, vala, Vala, vapi, vb.net, VB.net, vbnet, verilog, verilog, vhdl, VHDL, VHS, vhs, vim, VimL, vlang, vsh, vshell, vtt, vue, vue, vuejs, WDTE, WebGPU Shading Language, WebVTT, wgsl, Whiley, whiley, winbatch, XML, xml, xml+genshi, xml+kid, Xorg, xorg.conf, YAML, yaml, YANG, yang, z80, Z80 Assembly, Zed, zed, Zig, zig, zone, zsh
  - Styles: abap, algol, algol_nu, arduino, autumn, average, base16-snazzy, borland, bw, catppuccin-frappe, catppuccin-latte, catppuccin-macchiato, catppuccin-mocha, colorful, doom-one, doom-one2, dracula, emacs, evergarden, friendly, fruity, github, github-dark, gruvbox, gruvbox-light, hr_high_contrast, hrdark, igor, lovelace, manni, modus-operandi, modus-vivendi, monokai, monokailight, murphy, native, nord, nordic, onedark, onesenterprise, paraiso-dark, paraiso-light, pastie, perldoc, pygments, rainbow_dash, rose-pine, rose-pine-dawn, rose-pine-moon, rrt, solarized-dark, solarized-dark256, solarized-light, swapoff, tango, tokyonight-day, tokyonight-moon, tokyonight-night, tokyonight-storm, trac, vim, vs, vulcan, witchhazel, xcode, xcode-dark
  - Formatters: html, json, noop, svg, terminal, terminal16, terminal16m, terminal256, terminal8, tokens
 - openers:
    crng: get data from crypto/rand (accept a size limit as parameter)
    echo: echo the string given
    file: open local files
    gcs: get a GCP Cloud Storage object via gs://
    gemini: get URL via Gemini
    http: get URL via HTTP(S)
    curl: get URL via libcurl bindings
           libcurl/8.7.1 SecureTransport (LibreSSL/3.3.6) zlib/1.2.12 nghttp2/1.63.0
           protocols: dict,file,ftp,ftps,gopher,gophers,http,https,imap,imaps,ldap,ldaps,mqtt,pop3,pop3s,rtsp,smb,smbs,smtp,smtps,telnet,tftp
    mc: get a Minio-compatible object via mc:// (use ~/.mc/config.json or env for credentials)
    tcp: get data from listening on tcp://[HOST]:<PORT>
    prng: generate endless pcg rand (don't use for crypto) (accept a seed as parameter)
    s3: get an AWS s3 object via s3://
    ShellScp: get scp:// via local scp
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
        md5: compute the md5 checksum
        sha1: compute the sha1 checksum
        sha256: compute the sha256 checksum
        xxh32: compute the xxhash32 checksum
        xxh3: compute the xxh3 checksum
        xxh64: compute the xxhash64 checksum
    compress:
        bzip2: compress to bzip2 data (X:9 is compression level, 0-9)
        gzip: compress to gzip data (X:6 is compression level, 0-9)
        lz4: compress to lz4 data (X:0 is compression level, 0-9)
        lzma2: compress to lzma2 data
        lzma: compress to lzma data
        minlz: compress to minlz data
        pgzip: compress with pgzip  (X:6 is compression level, 0-9, blockSize, blocks)
        s2: compress to s2 data
        snap: compress to snappy data
        xz: compress to xz data
        zip: compress to zip data
        zlib: compress to zlib data (X:6 is compression level, 0-9)
        zstd: compress to zstd data (X:4 is compression level, 1-22)
    convert:
        base64: encode to base64
        feed2y: rss/atom/json feed -> YAML
        hex: dump in lowercase hex
        html2md: html -> markdown
        j2y: JSON -> YAML
        j5j: JSON5 -> JSON
        jcs: JSON -> JSON Canonicalization (RFC 8785)
        plist2Y: display an Apple plist as yaml
        qp: encode quoted-printable data
        unbase64: decode base64
        unhex: decode hex, ignore all non-hex chars
        unqp: decode quoted-printable data
        y2j: YAML -> JSON
    decompress:
        punbzip2: parallel decompress bzip2 data (X:0 is concurrency, 0 is auto)
        unbzip2: decompress bzip2 data
        ungzip: decompress gzip data
        unlz4: decompress lz4 data
        unlzfse: decompress lzfse data
        unlzma2: decompress lzma2 data
        unlzma: decompress lzma data
        unminlz: decompress minlz data
        unpgzip: decompress with pgzip
        uns2: decompress s2 data
        unsnap: decompress snappy data
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
    plugin:
        wasm: a wasi (wasm) module to apply (path as first argument
        yaegi: a yaegi script to apply (path as first argument, symbol as second argument)

  ('X:Y' means X is an argument with default value Y)

  mutator aliases:
    b64: base64
    cgpt: chatgpt
    d: discard
    dum, dumm: dummy
    gai: googleai
    xxd, hd: hexdump
    h2m, h2md: html2md
    hf: huggingface
    l, head: limit
    mime: mimetype
    mistral: mistralai
    unb64, ub64: unbase64
```
