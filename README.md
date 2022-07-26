# ccat
[![Go](https://github.com/batmac/ccat/actions/workflows/go.yml/badge.svg)](https://github.com/batmac/ccat/actions/workflows/go.yml)
![GitHub](https://img.shields.io/github/license/batmac/ccat)
[![Go Report Card](https://goreportcard.com/badge/github.com/batmac/ccat)](https://goreportcard.com/report/github.com/batmac/ccat)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbatmac%2Fccat.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbatmac%2Fccat?ref=badge_shield)
[![codecov](https://codecov.io/gh/batmac/ccat/branch/main/graph/badge.svg?token=PCD6DM6S75)](https://codecov.io/gh/batmac/ccat)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/batmac/ccat/main.svg)](https://results.pre-commit.ci/latest/github/batmac/ccat/main)

cat on steroids.
Leveraging great go modules to ease my CLI life.

## install
- get the latest release from https://github.com/batmac/ccat/releases/latest
- or `go install github.com/batmac/ccat/cmd/ccat@latest`

then, you can update to the latest github release with `ccat --selfupdate`

## build
you need a recent go, available build tags:
- `libcurl`: build with the libcurl opener.
- `fileonly`: build with the local file opener only.
- `nomd`: build without the markdown interpreter (glamour).
- `nohl`: build without the syntax-highlighter.
- `crappy`: build with some crappy (but useful) openers/mutators (needs a recent go version).

for instance:
`go build --tags libcurl,crappy .`

## examples
```
cat -X "kubectl get all" -i -w -t ready,Running --bg
NAME           READY   STATUS    RESTARTS   AGE
pod/busybox3   1/1     Running   1          17m

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.96.0.1    <none>        443/TCP   21d
```
( "READY and "Running" have different background colors )

```
$ echo Hello: | cat -m y2j,j2y,base64,hexdump
00000000  53 47 56 73 62 47 38 36  49 47 35 31 62 47 77 4b  |SGVsbG86IG51bGwK|
```
```
$ cat -m zstd LICENSE | cat -m base64
KLUv/WQwA+UUAGKyjCQQi1gA+Mu1uPH4trE1CdEhIptwP1oghhOOOP5rTRkCAKACAAdroEXGmA2Uz10inntYY97kIsI70zHae3CBb+DnAi4lAWNzxHPSRXalSceYEc4rc2FxPri8gX/QzYX9sT/32JzZHMbds2JzJmwIlmDMlYRjz/LGa8+BueftcBkb80Dn00G42FDA5RFeRujMBVvNmX8e5sjP1xLa4ftzH5ivJRriudg89+A789ibhIq4UDyyZ7G4GZni7BE2l7K96ZqlNKZ73IHNx5wltsQbeNlcc+bhPNOmS7jnE77PxeGI55jFszx3hWnuUSTcivBStqhQaFeY5xDeeOwNfGM+FiNjfC6gMx3jCoB0lwzgTGNcAoyKS/cWGWM2UL45c60B6ko+0c6WZuyo8Cl6249qdXyaqIRxSxnqfO0rw3XQV2mudYDt7dyRvdAhthvuNN3Hzc9QXbj4N6oSbktVj4U7cqWG33ycO0JVyW/+cW3JhdY0VXzqEkSWXpkcFb5Sc+VoKuOOrki4J0clqZErk0dHN3JFqbj6zQu/SrkjPUFsR9UE2b5ib2eoqYwr9VxdKXfDX6V822uiNRIAkACvUg44rlM/bm9nmm2NVNCNHtVMG7oY00eWvlK0IT8uPm4p/3HXRHKh00Ns26bZcW8vpqrkE0t/85UKoagE/XFfK/W4EtVsKQX5xFGpKh337W2MO0JVxVr2ZkuX14w7S+jhBgYFBheAAp4lSFMTyV3WqthWqJ9PlwgmAGYTiQo1Zi5cSUcnOBoEX3ev5YUMcVnvUrgx5HaCOLWhPcddUjOGZmX2yZFNAyB+ELmdEtsBPE17UsXU4fTV1jdNHpPl8KiAkC34tsZqWhxaa3RlIgn25jJkslDF+VBC8cW5mikplMUF0g==
```

## docker image
multi-arch container images with tags `libcurl,crappy` are automatically built by [Github Actions](https://github.com/batmac/ccat/actions/workflows/docker-images.yml) and made available on [Dockerhub](https://hub.docker.com/r/batmac/ccat) and [Github Packages](https://github.com/batmac/ccat/pkgs/container/ccat). <br/>
for instance:
```
docker run --rm ghcr.io/batmac/ccat:latest -h
```
or
```
kubectl run -i --tty ccat --image=batmac/ccat:latest -- /bin/sh
```

## help
```
version v0.9.19-12-g0ed004d [libcurl,crappy], commit 0ed004d76d74d47a604b451864e5222afd2341ea, built at 2022-07-26@06:06:39 by build.sh (go1.18.4 darwin/arm64)
  -t, --tokens string      comma-separated list of tokens
  -i, --ignore-case        tokens given with -t are case-insensitive
  -o, --only               don't display lines without at least one token
  -r, --raw                don't treat tokens as regexps
  -n, --line-number        number the output lines, starting at 1.
  -L, --flock-in           exclusively flock each file before reading
  -l, --flock-out          exclusively flock stdout
  -w, --word               read word by word instead of line by line (only works with utf8)
  -X, --exec string        command to exec on each file before processing it
  -b, --bg                 colorize the background instead of the font
  -H, --humanize           try to do what is needed to help (syntax-highlight, autodetect, etc.)
  -S, --style string       style to use (only used if -H, --fullhelp for the list)
  -F, --formatter string   formatter to use (only used if -H, --fullhelp for the list)
  -P, --lexer string       lexer to use (only used if -H, --fullhelp for the list)
  -m, --mutators string    mutators to use (comma-separated), --fullhelp for the list
  -V, --version            print version on stdout
      --license            print license on stdout
  -h, --help               print usage
      --fullhelp           print full usage
      --selfupdate         Update to latest Github release
      --check              Check version with the latest Github release
  -d, --debug              debug what we are doing
  -k, --insecure           get files insecurely (globally)

---
ccat <files>...
 - highlighter (used with -H):
  - Lexers: [1S 1S:Enterprise ABAP ABNF AL ANTLR APL ActionScript ActionScript 3 Ada Angular2 ApacheConf AppleScript Arduino ArmAsm Awk BNF Ballerina Base Makefile Bash BashSession Batchfile BibTeX Bicep BlitzBasic Brainfuck C C# C++ CFEngine3 CMake COBOL CSS Caddyfile Caddyfile Directives Cap'n Proto Cassandra CQL Ceylon ChaiScript Cheetah Clojure CoffeeScript Common Lisp Common Lisp Coq Crystal Cucumber Cython D DTD Dart Diff Django/Jinja Docker Dylan EBNF Elixir Elm EmacsLisp EmacsLisp Erlang FSharp Factor Fennel Fish Forth Fortran FortranFixed GAS GDScript GLSL Genshi Genshi HTML Genshi Text Gherkin Gherkin Gnuplot Go Go HTML Template Go HTML Template Go Text Template GraphQL Groff Groovy HCL HLB HTML HTTP Handlebars Haskell Haxe Hexdump Hy INI Idris Igor Io J JSON Java JavaScript Julia Jungle Kotlin LLVM Lighttpd configuration file Lua MLIR MZN Mako Mason Mathematica Matlab Meson Metal MiniZinc Modula-2 MonkeyC MorrowindScript MySQL Myghty NASM Newspeak Nginx configuration file Nim Nix OCaml Objective-C Octave OnesEnterprise OpenEdge ABL OpenSCAD Org Mode PHP PHTML PL/pgSQL POVRay PacmanConf Perl Pig PkgConfig Plutus Core Pony PostScript PostgreSQL SQL dialect PowerQuery PowerShell Prolog PromQL Protocol Buffer Puppet Python Python 2 QBasic QML R Racket Ragel Raku ReasonML Rexx Ruby Rust SAS SCSS SPARQL SQL SYSTEMD Sass Scala Scheme Scilab Sieve Smalltalk Smarty Snobol Solidity SquidConf Standard ML Stylus Svelte Swift TASM TOML TableGen Tcl Tcsh TeX Termcap Terminfo Terraform Thrift TradingView Transact-SQL Turing Turtle Twig TypeScript TypoScript TypoScriptCssData TypoScriptHtmlData V V shell VB.net VHDL VimL WDTE Whiley XML Xorg YAML YANG Zed Zig abap abl abnf aconf actionscript actionscript3 ada ada2005 ada95 al antlr apache apacheconf apl applescript arduino arexx armasm as as3 asm awk b3d ballerina bash bash-session basic bat batch bf bib bibtex bicep blitzbasic bnf bplus brainfuck bsdmake c c# c++ caddy caddy-d caddyfile caddyfile-d caddyfile-directives capnp cassandra ceylon cf3 cfengine3 cfg cfs cfstatement chai chaiscript cheetah cl cl clj clojure cmake cobol coffee coffee-script coffeescript common-lisp common-lisp console coq cpp cql cr crystal csh csharp css cucumber cython d dart diff django docker dockerfile dosbatch dosini dtd duby dylan ebnf elisp elisp elixir elm emacs emacs emacs-lisp emacs-lisp erlang ex exs factor fennel fish fishshell fnl forth fortran fortranfixed fsharp gas gawk gd gdscript genshi genshitext gherkin glsl gnuplot go go-html-template go-html-template go-text-template golang gql graphql graphqls groff groovy handlebars haskell haxe hbs hcl hexdump hlb hs html html+genshi html+kid http hx hxsl hylang idr idris igor igorpro ini io j java javascript jinja jl js json jsx julia jungle kid kotlin ksh latex lighttpd lighty lisp lisp llvm lua m2 make makefile mako man mariadb markdown mason mathematica matlab mawk mcfunction mcfunction md meson meson.build metal mf minizinc mkd mlir mma modula2 monkeyc morrowind mwscript myghty mysql mzn nasm nawk nb newspeak ng2 nginx nim nimrod nix nixos no-highlight nroff obj-c objc objective-c objectivec ocaml octave ones onesenterprise openedge openedgeabl openscad org orgmode pacmanconf perl perl6 php php3 php4 php5 phtml pig pkgconfig pl pl6 plain plaintext plc plpgsql plutus-core pony posh postgres postgresql postscr postscript pov powerquery powershell pq progress prolog promql proto protobuf ps1 psd1 psm1 puppet py py2 py3 pyrex python python2 python3 pyx qbasic qbs qml r racket ragel raku rb reStructuredText react react reason reasonml reg registry rest restructuredtext rexx rkt rs rst ruby rust s sage sas sass scala scheme scilab scm scss sh shell shell-session sieve smalltalk smarty sml snobol sol solidity sparql spitfire splus sql squeak squid squid.conf squidconf st stylus sv svelte swift systemd systemverilog systemverilog t-sql tablegen tasm tcl tcsh termcap terminfo terraform tex text tf thrift toml tradingview ts tsql tsx turing turtle tv twig typescript typoscript typoscriptcssdata typoscripthtmldata udiff v v vb.net vbnet verilog verilog vhdl vim vlang vsh vshell vue vue vuejs whiley winbatch xml xml+genshi xml+kid xorg.conf yaml yang zed zig zsh]
  - Styles: [abap algol algol_nu arduino autumn average base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github gruvbox hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarized-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark]
  - Formatters: [html json noop svg terminal terminal16 terminal16m terminal256 terminal8 tokens]
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
        md: Render Markdown (with glamour)
        mimetype: detect mimetype
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
