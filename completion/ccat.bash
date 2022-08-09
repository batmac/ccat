# ccat completion                                          -*- shell-script -*-

# This bash completions script was generated by
# completely (https://github.com/dannyben/completely)
# Modifying it manually is not recommended

_ccat_completions() {
  local cur=${COMP_WORDS[COMP_CWORD]}
  local compline="${COMP_WORDS[@]:1:$COMP_CWORD-1}"

  case "$compline" in
    ''*'--formatter')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "html json noop svg terminal terminal16 terminal16m terminal256 terminal8 token" -- "$cur" )
      ;;

    ''*'--mutators')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "sponge md5 cb dummy help hexdump indent j jcs mimetype translate wrap wrapU md5 sha1 sha256 xxh3 xxhash bzip2 gzip lz4 lzma2 lzma s2 snap xz zip zlib zstd base64 hex j2y plist2Y qp unbase64 unqp y2j unbzip2 ungzip unlz4 unlzfse unlzma2 unlzma uns2 unsnap unxz unzip unzlib unzstd filterUTF8 removeANSI" -- "$cur" )
      ;;

    ''*'--style')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "abap algol algol_nu arduino autumn average base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github gruvbox hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarize d-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark" -- "$cur" )
      ;;

    ''*'--lexer')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "1S 1S:Enterprise abap ABAP abl ABNF abnf aconf ActionScript actionscript actionscript3 ada Ada ada2005 ada95 al AL Angular2 antlr ANTLR apache apacheconf ApacheConf apl APL applescript AppleScript arduino Arduino arexx armasm ArmAsm as as3 asm Awk awk b3d ballerina Ballerina Base Makefile bash Bash bash-session BashSession basic bat batch Batchfile bf bib bibtex BibTeX bicep Bicep blitzbasic BlitzBasic bnf BNF bplus brainfuck Brainfuck bsdmake c C c# C# c++ C++ caddy caddy-d Caddyfile caddyfile Caddyfile Directives caddyfile-d caddyfile-directives capnp cassandra Cassandra CQL Ceylon ceylon cf3 CFEngine3 cfengine3 cfg cfs cfstatement chai chaiscript ChaiScript Cheetah cheetah cl cl clj clojure Clojure cmake CMake cobol COBOL coffee coffee-script coffeescript CoffeeScript common-lisp common-lisp console coq Coq cpp cql cr crystal Crystal csh csharp css CSS cucumber Cucumber Cython cython D d Dart dart Diff diff django Django/Jinja docker Docker dockerfile dosbatch dosini dtd DTD duby Dylan dylan EBNF ebnf elisp elisp elixir Elixir elm Elm emacs emacs emacs-lisp EmacsLisp erlang Erlang ex exs Factor factor Fennel fennel Fish fish fishshell fnl Forth forth Fortran fortran fortranfixed FortranFixed fsharp FSharp gas GAS gawk gd gdscript GDScript genshi Genshi Genshi HTML Genshi Text genshitext Gherkin gherkin Gherkin GLSL glsl gnuplot Gnuplot go Go Go HTML Template Go HTML Template Go Text Template go-html-template go-html-template go-text-template golang gql GraphQL graphql graphqls Groff groff groovy Groovy handlebars Handlebars Haskell haskell haxe Haxe hbs HCL hcl Hexdump hexdump HLB hlb hs html HTML html+genshi html+kid HTTP http hx hxsl Hy hylang idr Idris idris Igor igor igorpro ini INI Io io j J java Java javascript JavaScript jinja jl js JSON json jsx Julia julia Jungle jungle kid kotlin Kotlin ksh latex lighttpd Lighttpd configuration file lighty lisp lisp LLVM llvm Lua lua m2 make makefile Mako mako man mariadb markdown mason Mason mathematica Mathematica matlab Matlab mawk mcfunction mcfunction md Meson meson meson.build metal Metal mf MiniZinc minizinc mkd MLIR mlir mma Modula-2 modula2 monkeyc MonkeyC morrowind MorrowindScript mwscript myghty Myghty mysql MySQL MZN mzn NASM nasm nawk nb Newspeak newspeak ng2 nginx Nginx configuration file Nim nim nimrod nix Nix nixos no-highlight nroff obj-c objc objective-c Objective-C objectivec ocaml OCaml octave Octave ones onesenterprise OnesEnterprise openedge OpenEdge ABL openedgeabl openscad OpenSCAD org Org Mode orgmode PacmanConf pacmanconf Perl perl perl6 php PHP php3 php4 php5 PHTML phtml pig Pig PkgConfig pkgconfig pl PL/pgSQL pl6 plain plaintext plc plpgsql Plutus Core plutus-core pony Pony posh postgres postgresql PostgreSQL SQL dialect postscr postscript PostScript pov POVRay powerquery PowerQuery powershell PowerShell pq progress Prolog prolog PromQL promql proto protobuf Protocol Buffer ps1 psd1 psm1 puppet Puppet py py2 py3 pyrex python Python Python 2 python2 python3 pyx QBasic qbasic qbs qml QML R r Racket racket Ragel ragel raku Raku rb react react reason ReasonML reasonml reg registry rest reStructuredText restructuredtext Rexx rexx rkt rs rst ruby Ruby rust Rust s sage sas SAS Sass sass Scala scala Scheme scheme scilab Scilab scm scss SCSS sh shell shell-session sieve Sieve Smalltalk smalltalk Smarty smarty sml Snobol snobol sol Solidity solidity SPARQL sparql spitfire splus SQL sql squeak squid squid.conf squidconf SquidConf st Standard ML Stylus stylus sv svelte Svelte Swift swift SYSTEMD systemd systemverilog systemverilog t-sql tablegen TableGen tasm TASM tcl Tcl Tcsh tcsh Termcap termcap Terminfo terminfo Terraform terraform TeX tex text tf thrift Thrift toml TOML TradingView tradingview Transact-SQL ts tsql tsx turing Turing turtle Turtle tv twig Twig TypeScript typescript TypoScript typoscript typoscriptcssdata TypoScriptCssData TypoScriptHtmlData typoscripthtmldata udiff v v V V shell VB.net vb.net vbnet verilog verilog vhdl VHDL vim VimL vlang vsh vshell vue vue vuejs WDTE whiley Whiley winbatch xml XML xml+genshi xml+kid Xorg xorg.conf yaml YAML YANG yang Zed zed Zig zig zsh" -- "$cur" )
      ;;

    ''*'-P')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "1S 1S:Enterprise abap ABAP abl ABNF abnf aconf ActionScript actionscript actionscript3 ada Ada ada2005 ada95 al AL Angular2 antlr ANTLR apache apacheconf ApacheConf apl APL applescript AppleScript arduino Arduino arexx armasm ArmAsm as as3 asm Awk awk b3d ballerina Ballerina Base Makefile bash Bash bash-session BashSession basic bat batch Batchfile bf bib bibtex BibTeX bicep Bicep blitzbasic BlitzBasic bnf BNF bplus brainfuck Brainfuck bsdmake c C c# C# c++ C++ caddy caddy-d Caddyfile caddyfile Caddyfile Directives caddyfile-d caddyfile-directives capnp cassandra Cassandra CQL Ceylon ceylon cf3 CFEngine3 cfengine3 cfg cfs cfstatement chai chaiscript ChaiScript Cheetah cheetah cl cl clj clojure Clojure cmake CMake cobol COBOL coffee coffee-script coffeescript CoffeeScript common-lisp common-lisp console coq Coq cpp cql cr crystal Crystal csh csharp css CSS cucumber Cucumber Cython cython D d Dart dart Diff diff django Django/Jinja docker Docker dockerfile dosbatch dosini dtd DTD duby Dylan dylan EBNF ebnf elisp elisp elixir Elixir elm Elm emacs emacs emacs-lisp EmacsLisp erlang Erlang ex exs Factor factor Fennel fennel Fish fish fishshell fnl Forth forth Fortran fortran fortranfixed FortranFixed fsharp FSharp gas GAS gawk gd gdscript GDScript genshi Genshi Genshi HTML Genshi Text genshitext Gherkin gherkin Gherkin GLSL glsl gnuplot Gnuplot go Go Go HTML Template Go HTML Template Go Text Template go-html-template go-html-template go-text-template golang gql GraphQL graphql graphqls Groff groff groovy Groovy handlebars Handlebars Haskell haskell haxe Haxe hbs HCL hcl Hexdump hexdump HLB hlb hs html HTML html+genshi html+kid HTTP http hx hxsl Hy hylang idr Idris idris Igor igor igorpro ini INI Io io j J java Java javascript JavaScript jinja jl js JSON json jsx Julia julia Jungle jungle kid kotlin Kotlin ksh latex lighttpd Lighttpd configuration file lighty lisp lisp LLVM llvm Lua lua m2 make makefile Mako mako man mariadb markdown mason Mason mathematica Mathematica matlab Matlab mawk mcfunction mcfunction md Meson meson meson.build metal Metal mf MiniZinc minizinc mkd MLIR mlir mma Modula-2 modula2 monkeyc MonkeyC morrowind MorrowindScript mwscript myghty Myghty mysql MySQL MZN mzn NASM nasm nawk nb Newspeak newspeak ng2 nginx Nginx configuration file Nim nim nimrod nix Nix nixos no-highlight nroff obj-c objc objective-c Objective-C objectivec ocaml OCaml octave Octave ones onesenterprise OnesEnterprise openedge OpenEdge ABL openedgeabl openscad OpenSCAD org Org Mode orgmode PacmanConf pacmanconf Perl perl perl6 php PHP php3 php4 php5 PHTML phtml pig Pig PkgConfig pkgconfig pl PL/pgSQL pl6 plain plaintext plc plpgsql Plutus Core plutus-core pony Pony posh postgres postgresql PostgreSQL SQL dialect postscr postscript PostScript pov POVRay powerquery PowerQuery powershell PowerShell pq progress Prolog prolog PromQL promql proto protobuf Protocol Buffer ps1 psd1 psm1 puppet Puppet py py2 py3 pyrex python Python Python 2 python2 python3 pyx QBasic qbasic qbs qml QML R r Racket racket Ragel ragel raku Raku rb react react reason ReasonML reasonml reg registry rest reStructuredText restructuredtext Rexx rexx rkt rs rst ruby Ruby rust Rust s sage sas SAS Sass sass Scala scala Scheme scheme scilab Scilab scm scss SCSS sh shell shell-session sieve Sieve Smalltalk smalltalk Smarty smarty sml Snobol snobol sol Solidity solidity SPARQL sparql spitfire splus SQL sql squeak squid squid.conf squidconf SquidConf st Standard ML Stylus stylus sv svelte Svelte Swift swift SYSTEMD systemd systemverilog systemverilog t-sql tablegen TableGen tasm TASM tcl Tcl Tcsh tcsh Termcap termcap Terminfo terminfo Terraform terraform TeX tex text tf thrift Thrift toml TOML TradingView tradingview Transact-SQL ts tsql tsx turing Turing turtle Turtle tv twig Twig TypeScript typescript TypoScript typoscript typoscriptcssdata TypoScriptCssData TypoScriptHtmlData typoscripthtmldata udiff v v V V shell VB.net vb.net vbnet verilog verilog vhdl VHDL vim VimL vlang vsh vshell vue vue vuejs WDTE whiley Whiley winbatch xml XML xml+genshi xml+kid Xorg xorg.conf yaml YAML YANG yang Zed zed Zig zig zsh" -- "$cur" )
      ;;

    ''*'-m')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "sponge md5 cb dummy help hexdump indent j jcs mimetype translate wrap wrapU md5 sha1 sha256 xxh3 xxhash bzip2 gzip lz4 lzma2 lzma s2 snap xz zip zlib zstd base64 hex j2y plist2Y qp unbase64 unqp y2j unbzip2 ungzip unlz4 unlzfse unlzma2 unlzma uns2 unsnap unxz unzip unzlib unzstd filterUTF8 removeANSI" -- "$cur" )
      ;;

    ''*'-F')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "html json noop svg terminal terminal16 terminal16m terminal256 terminal8 token" -- "$cur" )
      ;;

    ''*'-S')
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -W "abap algol algol_nu arduino autumn average base16-snazzy borland bw colorful doom-one doom-one2 dracula emacs friendly fruity github gruvbox hr_high_contrast hrdark igor lovelace manni monokai monokailight murphy native nord onesenterprise paraiso-dark paraiso-light pastie perldoc pygments rainbow_dash rrt solarized-dark solarized-dark256 solarize d-light swapoff tango trac vim vs vulcan witchhazel xcode xcode-dark" -- "$cur" )
      ;;

    *)
      while read; do COMPREPLY+=( "$REPLY" ); done < <( compgen -A file -W "http:// https:// s3:// gcs:// tcp:// mc:// file:// --version --help -t --tokens -i --ignore-case -o --only -r --raw -n --line-number -L --flock-in -l --flock-out -w --word -X --exec string -b --bg -H --humanize -S --style string -F --formatter string -P --lexer string -m --mutators string -V --version --license --gomod -h --help --fullhelp --selfupdate --check -d --debug -k --insecure" -- "$cur" )
      ;;

  esac
} &&
complete -F _ccat_completions ccat

# ex: filetype=sh
