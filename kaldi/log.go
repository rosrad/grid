// log config for kaldi task
package kaldi

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func LogFile() string {
	log_dir := path.Join(RootPath(), "log")
	InsureDir(log_dir)
	return path.Join(log_dir, Now()) + ".log"
}

var g_logfile *os.File
var g_warn, g_err, g_trace *log.Logger

func Init() {
	var err error
	g_logfile, err = os.OpenFile(LogFile(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		g_logfile = nil
		log.Println("Error of log file:", err)
	}
	Warn().Println("")
	Warn().Println("=============================")
	Warn().Println("")
}

func Uninit() {
	Warn().Println("")
	Warn().Println("=============================")
	Warn().Println("")
	if g_logfile != nil {
		defer g_logfile.Close()
	}
}

func logger(tag string) *log.Logger {
	multi := io.MultiWriter(os.Stdout, g_logfile)
	return log.New(multi,
		strings.ToUpper(tag)+" :",
		log.Lshortfile)
}

func Err() *log.Logger {
	if g_err == nil {
		g_err = logger("Error")
	}
	return g_err
}

func Warn() *log.Logger {
	if g_warn == nil {
		g_warn = logger("Warn")
	}
	return g_warn
}

func Trace() *log.Logger {
	if g_trace == nil {
		g_trace = logger("Trace")
	}
	return g_trace

}
