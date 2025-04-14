package mutators

import (
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"

	"github.com/batmac/ccat/pkg/globalctx"
	"github.com/batmac/ccat/pkg/log"
	"github.com/batmac/ccat/pkg/term"
)

func init() {
	singleRegister("it2dl", it2dl,
		withDescription("download via iTerm2 escape code, does not work in other terminals. Must be the last mutator of the pipeline."),
		withCategory("convert"),
		withConfigBuilder(stdConfigStringWithDefault("")),
		withExpectingBinary(), // don't highlight the output as it must unmodified
		withExpectingFinal(),  // this mutator must be the last in the pipeline
	)
}

func it2dl(w io.WriteCloser, r io.ReadCloser, config any) (int64, error) {
	if !term.IsITerm2() {
		log.Println("WARNING: it2dl is designed for iTerm2. It may not work as expected in this terminals.")
	} else if !term.IsStdoutTerminal() {
		log.Println("WARNING: it2dl is designed for iTerm2. It may not work as expected in non-terminal outputs.")
	}

	// here is the shell version from upstream: (notably, the size must be known in advance)
	// function b64_encode() {
	// 	load_version
	// 	if [[ "$IT2DL_BASE64_VERSION" =~ GNU ]]; then
	// 		# Disable line wrap
	// 		base64 -w0
	// 	else
	// 		base64
	// 	fi
	// }
	// for fn in "$@"
	// do
	//   if [ -r "$fn" ] ; then
	// 	[ -d "$fn" ] && { echo "$fn is a directory"; continue; }
	// 	printf '\033]1337;File=name=%s;' $(echo -n "$fn" | b64_encode)
	// 	wc -c "$fn" | awk '{printf "size=%d",$1}'
	// 	printf ":"
	// 	base64 < "$fn"
	// 	printf '\a'
	//   else
	// 	echo File $fn does not exist or is not readable.
	//   fi
	// done

	// name := filepath.Base(globalctx.Get("path").(string))
	// if name == "" {
	// fileName := base64.StdEncoding.EncodeToString([]byte(config.(string)))
	fileName, ok := config.(string)
	if !ok || fileName == "" {
		fileName, ok = globalctx.Get("path").(string)
		if !ok || fileName == "" {
			fileName = "file.bin"
		} else {
			fileName = filepath.Base(fileName)
		}
	}
	log.Debugf("fileName: %s", fileName)

	fileName = base64.StdEncoding.EncodeToString([]byte(fileName))

	data, err := io.ReadAll(r)
	if err != nil || len(data) == 0 {
		return 0, err
	}

	preWritten, err := io.WriteString(w,
		fmt.Sprintf("\033]1337;File=name=%s;size=%d:",
			fileName, len(data),
		))
	if err != nil {
		return 0, err
	}

	// base64 encoder
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	// write the data to the encoder
	_, err = encoder.Write(data)
	if err != nil {
		return 0, err
	}
	// close the encoder
	err = encoder.Close()
	if err != nil {
		return 0, err
	}
	// write the escape sequence to the writer
	_, err = io.WriteString(w, "\a")
	if err != nil {
		return 0, err
	}
	// // close the writer
	// err = w.Close()
	// if err != nil {
	// 	return 0, err
	// }

	return int64(preWritten), nil
}
