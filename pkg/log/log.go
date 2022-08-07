package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
)

var (
	flags  = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	Debug  = &Logger{Logger: log.New(io.Discard, "", flags)}
	Stderr = &Logger{Logger: log.New(os.Stderr, "", flags)}

	DebugIsDiscard  int32
	continueOnFatal int32
)

func init() {
}

type Logger struct {
	*log.Logger
}

func Default() *Logger {
	return Debug
}

func SetDebug(w io.Writer) {
	var isDiscard int32 = 0
	if w == io.Discard {
		isDiscard = 1
	}
	atomic.StoreInt32(&DebugIsDiscard, isDiscard)
	Debug.SetOutput(w)
}

func Debugf(format string, v ...interface{}) {
	if atomic.LoadInt32(&DebugIsDiscard) != 0 {
		return
	}
	_ = Debug.Output(2, fmt.Sprintf(format, v...))
}

func Debugln(v ...interface{}) {
	if atomic.LoadInt32(&DebugIsDiscard) != 0 {
		return
	}
	_ = Debug.Output(2, fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	Stderr.Output(2, fmt.Sprintf(format, v...))
}

func Print(v ...interface{}) {
	Stderr.Output(2, fmt.Sprintln(v...))
}

func Println(v ...interface{}) {
	Stderr.Output(2, fmt.Sprint(v...))
}

func Fatal(v ...interface{}) {
	Stderr.Output(2, fmt.Sprint(v...))
	if atomic.LoadInt32(&continueOnFatal) <= 0 {
		os.Exit(1)
	}
}

// pretty print stuff
func Pp(data interface{}) string {
	var j []byte
	//    var err := error
	j, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		Stderr.Println(err)
		return err.Error()
	}
	return string(j)
}

func SetContinueOnFatal() {
	atomic.StoreInt32(&continueOnFatal, 1)
}
