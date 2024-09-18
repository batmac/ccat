#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>

// zig cc --target=wasm32-wasi main.c -o c.wasm

int main(int argc, char **argv)
{
    ssize_t n, m;
    char buf[BUFSIZ];

    int in = STDIN_FILENO;
    int out = STDOUT_FILENO;

    while ((n = read(in, buf, BUFSIZ)) > 0)
    {
        char *ptr = buf;
        while (n > 0)
        {
            m = write(out, ptr, (size_t)n);
            if (m < 0)
            {
                fprintf(stderr, "write error: %s\n", strerror(errno));
                exit(1);
            }
            n -= m;
            ptr += m;
        }
    }

    if (n < 0)
    {
        fprintf(stderr, "read error: %s\n", strerror(errno));
        exit(1);
    }

    return EXIT_SUCCESS;
}
