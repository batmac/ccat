package log

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var (
	flags  = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	Debug  = &Logger{Logger: log.New(os.Stderr, "", flags)}
	Stderr = &Logger{Logger: log.New(os.Stderr, "", flags)}
)

type Logger struct {
	*log.Logger
}

func Default() *Logger {
	return Debug
}

func SetDebug(w io.Writer) {
	Debug.SetOutput(w)
}
func Debugf(format string, v ...interface{}) {
	Debug.Printf(format, v...)
}
func Debugln(v ...interface{}) {
	Debug.Println(v...)
}

func Printf(format string, v ...interface{}) {
	Stderr.Printf(format, v...)
}
func Print(v ...interface{}) {
	Stderr.Print(v...)
}
func Println(v ...interface{}) {
	Stderr.Println(v...)
}

func Fatal(v ...interface{}) {
	Stderr.Fatal(v...)
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
