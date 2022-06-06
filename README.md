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
  -i    case-insensitive
  -l    exclusively flock stdout
  -n    number the output lines, starting at 1.
  -o    don't display lines without at least one token
  -r    don't treat tokens as regexps
  -t string
        comma-separated list of tokens
  -w    read word by word instead of line by line (only works with utf8)
```
