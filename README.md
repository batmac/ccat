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
  buildAndTest*      buildDefault,test
  buildDefault       tags: libcurl,crappy
  buildMinimal       tags: nohl,fileonly
  clean
  install            put ccat to $GOPATH/bin/ccat
  installDeps        go mod download
  test               all
  testCompression    test_compression_e2e
  testGo             go test ./...
  updateREADME
  verifyDeps         go mod verify

* default target
$ go run magefiles/mage.go # or 'mage'
```

## Update

- You can update to the latest github release with `ccat --selfupdate`.
- You can check against your current installed version with `ccat --check`.

## Build tags

available build tags:

- `libcurl`: build with the libcurl opener.
- `fileonly`: build with the local file opener only.
- `nohl`: build without the syntax-highlighter.
- `crappy`: build with some crappy (but useful) openers/mutators.

## Examples

```shell
$ cat -X "kubectl get all" -i -w -t ready,Running --bg
NAME           READY   STATUS    RESTARTS   AGE
pod/busybox3   1/1     Running   1          17m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   21d
```

( "READY and "Running" have different background colors )

```shell
$ echo Hello: | cat -m y2j,j2y,base64,hexdump
00000000  53 47 56 73 62 47 38 36  49 47 35 31 62 47 77 4b  |SGVsbG86IG51bGwK|
```

```shell
$ cat -m zstd LICENSE | cat -m base64
KLUv/WQwA+UUAGKyjCQQi1gA+Mu1uPH4trE1CdEhIptwP1oghhOOOP5rTRkCAKACAAdroEXGmA2Uz10inntYY97kIsI70zHae3CBb+DnAi4lAWNzxHPSRXalSceYEc4rc2FxPri8gX/QzYX9sT/32JzZHMbds2JzJmwIlmDMlYRjz/LGa8+BueftcBkb80Dn00G42FDA5RFeRujMBVvNmX8e5sjP1xLa4ftzH5ivJRriudg89+A789ibhIq4UDyyZ7G4GZni7BE2l7K96ZqlNKZ73IHNx5wltsQbeNlcc+bhPNOmS7jnE77PxeGI55jFszx3hWnuUSTcivBStqhQaFeY5xDeeOwNfGM+FiNjfC6gMx3jCoB0lwzgTGNcAoyKS/cWGWM2UL45c60B6ko+0c6WZuyo8Cl6249qdXyaqIRxSxnqfO0rw3XQV2mudYDt7dyRvdAhthvuNN3Hzc9QXbj4N6oSbktVj4U7cqWG33ycO0JVyW/+cW3JhdY0VXzqEkSWXpkcFb5Sc+VoKuOOrki4J0clqZErk0dHN3JFqbj6zQu/SrkjPUFsR9UE2b5ib2eoqYwr9VxdKXfDX6V822uiNRIAkACvUg44rlM/bm9nmm2NVNCNHtVMG7oY00eWvlK0IT8uPm4p/3HXRHKh00Ns26bZcW8vpqrkE0t/85UKoagE/XFfK/W4EtVsKQX5xFGpKh337W2MO0JVxVr2ZkuX14w7S+jhBgYFBheAAp4lSFMTyV3WqthWqJ9PlwgmAGYTiQo1Zi5cSUcnOBoEX3ev5YUMcVnvUrgx5HaCOLWhPcddUjOGZmX2yZFNAyB+ELmdEtsBPE17UsXU4fTV1jdNHpPl8KiAkC34tsZqWhxaa3RlIgn25jJkslDF+VBC8cW5mikplMUF0g==
```

## Docker image

multi-arch container images with tags `libcurl,crappy` are automatically built by [Github Actions](https://github.com/batmac/ccat/actions/workflows/docker-images.yml) and made available on [Dockerhub](https://hub.docker.com/r/batmac/ccat) and [Github Packages](https://github.com/batmac/ccat/pkgs/container/ccat). <br/>
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
version v0.9.23-16-g46d918b [libcurl,crappy], commit 46d918b44f8271a49e71157c6054c861b183e19f, built at 2022-08-16@22:26:17+0200 by Mage (go1.19 darwin/arm64)
  -t, --tokens string       comma-separated list of tokens
  -i, --ignore-case         tokens given with -t are case-insensitive
  -o, --only                don't display lines without at least one token
  -r, --raw                 don't treat tokens as regexps
  -n, --line-number         number the output lines, starting at 1.
  -L, --flock-in            exclusively flock each file before reading
  -l, --flock-out           exclusively flock stdout
  -w, --word                read word by word instead of line by line (only works with utf8)
  -X, --exec string         command to exec on each file before processing it
  -b, --bg                  colorize the background instead of the font
  -H, --humanize            try to do what is needed to help (syntax-highlight, autodetect, etc.)
  -S, --style string        style to use (only used if -H, --fullhelp for the list)
  -F, --formatter string    formatter to use (only used if -H, --fullhelp for the list)
  -P, --lexer string        lexer to use (only used if -H, --fullhelp for the list)
  -m, --mutators string     mutators to use (comma-separated), --fullhelp for the list
  -V, --version             print version on stdout
      --license             print license on stdout
  -B, --buildinfo           print build info on stdout
  -h, --help                print usage
      --fullhelp            print full usage
      --selfupdate          Update to latest Github release
      --check               Check version with the latest Github release
  -d, --debug               debug what we are doing
  -k, --insecure            get files insecurely (globally)
  -C, --completion string   print shell completion script
  -T, --ui                  display with a minimal ui

---
ccat <files>...
 - highlighter (used with -H):
  - Lexers: 1S, 1S:Enterprise, abap, ABAP, abl, ABNF, abnf, aconf, ActionScript, actionscript, ActionScript 3, actionscript3, ada, Ada, ada2005, ada95, al, AL, Angular2, antlr, ANTLR, apache, apacheconf, ApacheConf, apl, APL, applescript, AppleScript, arduino, Arduino, arexx, armasm, ArmAsm, as, as3, asm, Awk, awk, b3d, ballerina, Ballerina, Base Makefile, bash, Bash, bash-session, BashSession, basic, bat, batch, Batchfile, bf, bib, bibtex, BibTeX, bicep, Bicep, blitzbasic, BlitzBasic, bnf, BNF, bplus, brainfuck, Brainfuck, bsdmake, c, C, c#, C#, c++, C++, caddy, caddy-d, Caddyfile, caddyfile, Caddyfile Directives, caddyfile-d, caddyfile-directives, Cap'n Proto, capnp, cassandra, Cassandra CQL, Ceylon, ceylon, cf3, CFEngine3, cfengine3, cfg, cfs, cfstatement, chai, chaiscript, ChaiScript, Cheetah, cheetah, cl, cl, clj, clojure, Clojure, cmake, CMake, cobol, COBOL, coffee, coffee-script, coffeescript, CoffeeScript, Common Lisp, Common Lisp, common-lisp, common-lisp, console, coq, Coq, cpp, cql, cr, crystal, Crystal, csh, csharp, css, CSS, cucumber, Cucumber, Cython, cython, D, d, Dart, dart, Diff, diff, django, Django/Jinja, docker, Docker, dockerfile, dosbatch, dosini, dtd, DTD, duby, Dylan, dylan, EBNF, ebnf, elisp, elisp, elixir, Elixir, elm, Elm, emacs, emacs, emacs-lisp, emacs-lisp, EmacsLisp, EmacsLisp, erlang, Erlang, ex, exs, Factor, factor, Fennel, fennel, Fish, fish, fishshell, fnl, Forth, forth, Fortran, fortran, fortranfixed, FortranFixed, fsharp, FSharp, gas, GAS, gawk, gd, gdscript, GDScript, genshi, Genshi, Genshi HTML, Genshi Text, genshitext, Gherkin, gherkin, Gherkin, GLSL, glsl, gnuplot, Gnuplot, go, Go, Go HTML Template, Go HTML Template, Go Text Template, go-html-template, go-html-template, go-text-template, golang, gql, GraphQL, graphql, graphqls, Groff, groff, groovy, Groovy, handlebars, Handlebars, Haskell, haskell, haxe, Haxe, hbs, HCL, hcl, Hexdump, hexdump, HLB, hlb, hs, html, HTML, html+genshi, html+kid, HTTP, http, hx, hxsl, Hy, hylang, idr, Idris, idris, Igor, igor, igorpro, ini, INI, Io, io, j, J, java, Java, javascript, JavaScript, jinja, jl, js, JSON, json, jsx, Julia, julia, Jungle, jungle, kid, kotlin, Kotlin, ksh, latex, lighttpd, Lighttpd configuration file, lighty, lisp, lisp, LLVM, llvm, Lua, lua, m2, make, makefile, Mako, mako, man, mariadb, markdown, mason, Mason, mathematica, Mathematica, matlab, Matlab, mawk, mcfunction, mcfunction, md, Meson, meson, meson.build, metal, Metal, mf, MiniZinc, minizinc, mkd, MLIR, mlir, mma, Modula-2, modula2, monkeyc, MonkeyC, morrowind, MorrowindScript, mwscript, myghty, Myghty, mysql, MySQL, MZN, mzn, NASM, nasm, nawk, nb, Newspeak, newspeak, ng2, nginx, Nginx configuration file, Nim, nim, nimrod, nix, Nix, nixos, no-highlight, nroff, obj-c, objc, objective-c, Objective-C, objectivec, ocaml, OCaml, octave, Octave, ones, onesenterprise, OnesEnterprise, openedge, OpenEdge ABL, openedgeabl, openscad, OpenSCAD, org, Org Mode, orgmode, PacmanConf, pacmanconf, Perl, perl, perl6, php, PHP, php3, php4, php5, PHTML, phtml, pig, Pig, PkgConfig, pkgconfig, pl, PL/pgSQL, pl6, plain, plaintext, plc, plpgsql, Plutus Core, plutus-core, pony, Pony, posh, postgres, postgresql, PostgreSQL SQL dialect, postscr, postscript, PostScript, pov, POVRay, powerquery, PowerQuery, powershell, PowerShell, pq, progress, Prolog, prolog, PromQL, promql, proto, protobuf, Protocol Buffer, ps1, psd1, psm1, puppet, Puppet, py, py2, py3, pyrex, python, Python, Python 2, python2, python3, pyx, QBasic, qbasic, qbs, qml, QML, R, r, Racket, racket, Ragel, ragel, raku, Raku, rb, react, react, reason, ReasonML, reasonml, reg, registry, rest, reStructuredText, restructuredtext, Rexx, rexx, rkt, rs, rst, ruby, Ruby, rust, Rust, s, sage, sas, SAS, Sass, sass, Scala, scala, Scheme, scheme, scilab, Scilab, scm, scss, SCSS, sh, shell, shell-session, sieve, Sieve, Smalltalk, smalltalk, Smarty, smarty, sml, Snobol, snobol, sol, Solidity, solidity, SPARQL, sparql, spitfire, splus, SQL, sql, squeak, squid, squid.conf, squidconf, SquidConf, st, Standard ML, Stylus, stylus, sv, svelte, Svelte, Swift, swift, SYSTEMD, systemd, systemverilog, systemverilog, t-sql, tablegen, TableGen, tasm, TASM, tcl, Tcl, Tcsh, tcsh, Termcap, termcap, Terminfo, terminfo, Terraform, terraform, TeX, tex, text, tf, thrift, Thrift, toml, TOML, TradingView, tradingview, Transact-SQL, ts, tsql, tsx, turing, Turing, turtle, Turtle, tv, twig, Twig, TypeScript, typescript, TypoScript, typoscript, typoscriptcssdata, TypoScriptCssData, TypoScriptHtmlData, typoscripthtmldata, udiff, v, v, V, V shell, VB.net, vb.net, vbnet, verilog, verilog, vhdl, VHDL, vim, VimL, vlang, vsh, vshell, vue, vue, vuejs, WDTE, whiley, Whiley, winbatch, xml, XML, xml+genshi, xml+kid, Xorg, xorg.conf, yaml, YAML, YANG, yang, Zed, zed, Zig, zig, zsh
  - Styles: abap, algol, algol_nu, arduino, autumn, average, base16-snazzy, borland, bw, colorful, doom-one, doom-one2, dracula, emacs, friendly, fruity, github, gruvbox, hr_high_contrast, hrdark, igor, lovelace, manni, monokai, monokailight, murphy, native, nord, onesenterprise, paraiso-dark, paraiso-light, pastie, perldoc, pygments, rainbow_dash, rrt, solarized-dark, solarized-dark256, solarized-light, swapoff, tango, trac, vim, vs, vulcan, witchhazel, xcode, xcode-dark
  - Formatters: html, json, noop, svg, terminal, terminal16, terminal16m, terminal256, terminal8, tokens
 - openers:
    file: open local files
    gcs: get a GCP Cloud Storage object via gs://
    curl: get URL via libcurl bindings
           libcurl/7.79.1 SecureTransport (LibreSSL/3.3.6) zlib/1.2.11 nghttp2/1.45.1
           protocols: dict,file,ftp,ftps,gopher,gophers,http,https,imap,imaps,ldap,ldaps,mqtt,pop3,pop3s,rtsp,smb,smbs,smtp,smtps,telnet,tftp
    mc: get a Minio-compatible object via mc:// (use ~/.mc/config.json or env for credentials)
    tcp: get data from listening on tcp://[HOST]:<PORT>
    s3: get an AWS s3 object via s3://
    ShellScp: get scp:// via local scp

 - mutators:
        cb: put a copy in the clipboard
        dummy: a simple fifo
        help: display mutators help
        hexdump: dump in hex as xxd
        indent: indent the text with 4 chars
        j: JSON Re-indent
        jcs: JSON Canonicalization (RFC 8785)
        mimetype: detect mimetype
        sponge: soak all input before outputting it.
        translate: translate to $TARGET_LANGUAGE with google translate (need a valid key in $GOOGLE_API_KEY)
        wrap: word-wrap the text to 80 chars maximum
        wrapU: unconditionally wrap the text to 80 chars maximum
    checksum:
        md5: compute the md5 checksum
        sha1: compute the sha1 checksum
        sha256: compute the sha256 checksum
        xxh3: compute the xxh3 checksum
        xxhash: compute the xxhash (xxh64) checksum
    compress:
        bzip2: compress to bzip2 data
        gzip: compress to gzip data
        lz4: compress to lz4 data
        lzma2: compress to lzma2 data
        lzma: compress to lzma data
        s2: compress to s2 data
        snap: compress to snappy data
        xz: compress to xz data
        zip: compress to zip data
        zlib: compress to zlib data
        zstd: compress to zstd data
    convert:
        base64: encode to base64
        hex: dump in lowercase hex
        j2y: JSON -> YAML
        plist2Y: display an Apple plist as yaml
        qp: encode quoted-printable data
        unbase64: decode base64
        unqp: decode quoted-printable data
        y2j: YAML -> JSON
    decompress:
        unbzip2: decompress bzip2 data
        ungzip: decompress gzip data
        unlz4: decompress lz4 data
        unlzfse: decompress lzfse data
        unlzma2: decompress lzma2 data
        unlzma: decompress lzma data
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
    filter:
        filterUTF8: remove non-utf8
        removeANSI: remove ANSI codes

```
