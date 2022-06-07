# ccat
cat on steroids (mainly useful to multi-color lines/words according to tokens)

```
Usage of ccat:
 
  -F string
        formatter to use (only used if -H, look in -d for the list)
  -H    try to syntax-highlight
  -L    exclusively flock each file before reading
  -P string
        lexer to use (only used if -H, look in -d for the list)
  -S string
        style to use (only used if -H, look in -d for the list)
  -X string
        command to exec on each file before processing it
  -bg
        colorize the background instead of the font
  -d    debug what we are doing
  -i    tokens given with -t are case-insensitive
  -l    exclusively flock stdout
  -n    number the output lines, starting at 1.
  -o    don't display lines without at least one token
  -r    don't treat tokens as regexps
  -t string
        comma-separated list of tokens
  -w    read word by word instead of line by line (only works with utf8)
---
ccat <files>...
 - openers: [curl file http]
 - highlighter (-H):
  - Lexers: [1S 1S:Enterprise ABAP ABNF AL ANTLR APL ActionScript ActionScript 3 Ada Angular2 ApacheConf AppleScript Arduino ArmAsm Awk BNF Ballerina Base Makefile Bash BashSession Batchfile BibTeX Bicep BlitzBasic Brainfuck C C# C++ CFEngine3 CMake COBOL CSS Caddyfile Caddyfile Directives Cap'n Proto Cassandra CQL Ceylon ChaiScript Cheetah Clojure CoffeeScript Common Lisp Coq Crystal Cucumber Cython D DTD Dart Diff Django/Jinja Docker Dylan EBNF Elixir Elm EmacsLisp Erlang FSharp Factor Fennel Fish Forth Fortran FortranFixed GAS GDScript GLSL Genshi Genshi HTML Genshi Text Gherkin Gherkin Gnuplot Go Go HTML Template Go Text Template GraphQL Groff Groovy HCL HLB HTML HTTP Handlebars Haskell Haxe Hexdump Hy INI Idris Igor Io J JSON Java JavaScript Julia Jungle Kotlin LLVM Lighttpd configuration file Lua MLIR MZN Mako Mason Mathematica Matlab Meson Metal MiniZinc Modula-2 MonkeyC MorrowindScript MySQL Myghty NASM Newspeak Nginx configuration file Nim Nix OCaml Objective-C Octave OnesEnterprise OpenEdge ABL OpenSCAD Org Mode PHP PHTML PL/pgSQL POVRay PacmanConf Perl Pig PkgConfig Plutus Core Pony PostScript PostgreSQL SQL dialect PowerQuery PowerShell Prolog PromQL Protocol Buffer Puppet Python Python 2 QBasic QML R Racket Ragel Raku ReasonML Rexx Ruby Rust SAS SCSS SPARQL SQL SYSTEMD Sass Scala Scheme Scilab Sieve Smalltalk Smarty Snobol Solidity SquidConf Standard ML Stylus Svelte Swift TASM TOML TableGen Tcl Tcsh TeX Termcap Terminfo Terraform Thrift TradingView Transact-SQL Turing Turtle Twig TypeScript TypoScript TypoScriptCssData TypoScriptHtmlData VB.net VHDL VimL WDTE XML Xorg YAML YANG Zed Zig abap abl abnf aconf actionscript actionscript3 ada ada2005 ada95 al antlr apache apacheconf apl applescript arduino arexx armasm as as3 asm awk b3d ballerina bash bash-session basic bat batch bf bib bibtex bicep blitzbasic bnf bplus brainfuck bsdmake c c# c++ caddy caddy-d caddyfile caddyfile-d caddyfile-directives capnp cassandra ceylon cf3 cfengine3 cfg cfs cfstatement chai chaiscript cheetah cl clj clojure cmake cobol coffee coffee-script coffeescript common-lisp console coq cpp cql cr crystal csh csharp css cucumber cython d dart diff django docker dockerfile dosbatch dosini dtd duby dylan ebnf elisp elixir elm emacs emacs-lisp erlang ex exs factor fennel fish fishshell fnl forth fortran fortranfixed fsharp gas gawk gd gdscript genshi genshitext gherkin glsl gnuplot go go-html-template go-text-template golang gql graphql graphqls groff groovy handlebars haskell haxe hbs hcl hexdump hlb hs html html+genshi html+kid http hx hxsl hylang idr idris igor igorpro ini io j java javascript jinja jl js json jsx julia jungle kid kotlin ksh latex lighttpd lighty lisp llvm lua m2 make makefile mako man markdown mason mathematica matlab mawk mcfunction mcfunction md meson meson.build metal mf minizinc mkd mlir mma modula2 monkeyc morrowind mwscript myghty mysql mzn nasm nawk nb newspeak ng2 nginx nim nimrod nix nixos no-highlight nroff obj-c objc objective-c objectivec ocaml octave ones onesenterprise openedge openedgeabl openscad org orgmode pacmanconf perl perl6 php php3 php4 php5 phtml pig pkgconfig pl pl6 plain plaintext plc plpgsql plutus-core pony posh postgres postgresql postscr postscript pov powerquery powershell pq progress prolog promql proto protobuf ps1 psd1 psm1 puppet py py2 py3 pyrex python python2 python3 pyx qbasic qbs qml r racket ragel raku rb reStructuredText react react reason reasonml reg registry rest restructuredtext rexx rkt rs rst ruby rust s sage sas sass scala scheme scilab scm scss sh shell shell-session sieve smalltalk smarty sml snobol sol solidity sparql spitfire splus sql squeak squid squid.conf squidconf st stylus sv svelte swift systemd systemverilog systemverilog t-sql tablegen tasm tcl tcsh termcap terminfo terraform tex text tf thrift toml tradingview ts tsql tsx turing turtle tv twig typescript typoscript typoscriptcssdata typoscripthtmldata udiff v vb.net vbnet verilog verilog vhdl vim vue vue vuejs winbatch xml xml+genshi xml+kid xorg.conf yaml yang zed zig zsh]
  - Styles: [abap algol algol_nu arduino autumn base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarized-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark]
  - Formatters: [html json noop svg terminal terminal16 terminal16m terminal256 terminal8 tokens]
```
