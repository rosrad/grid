// log config for kaldi task
package kaldi

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func LogFile(tag string) string {
	dir := fmt.Sprintf("%d-%02d", time.Now().Year(), time.Now().Month())
	filename := fmt.Sprintf("%02d", time.Now().Day())
	log_dir := path.Join(RootPath(), "log", tag, dir)
	InsureDir(log_dir)
	return path.Join(log_dir, filename) + ".log"
}

var g_logfile *os.File
var g_warn, g_err, g_trace *log.Logger

func Init(root, LM string) {
	LoadGlobalConf()
	if "" != root {
		SysConf().Root = root
	}
	if "" != LM {
		SysConf().LM = LM
	}
	Trace().Println("LM :", SysConf().LM)
	Trace().Println("Root :", RootPath())
	LoadDataConf()
	var err error
	g_logfile, err = os.OpenFile(LogFile("default"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		g_logfile = nil
		log.Println("Error of log file:", err)
	}
	Warn().Println("")
	Warn().Println("=============================")
	Warn().Println("Default LogFile:", LogFile("default"))
	Warn().Println("=============================")
	Warn().Println("")
}

func NewLogWriter(tag string) io.Writer {
	if tag == "" {
		tag = "default"
	}
	f, err := os.OpenFile(LogFile(tag), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error of log file:", err)
	}
	return io.MultiWriter(os.Stdout, f)
}

func LogWriter() io.Writer {
	return io.MultiWriter(os.Stdout, g_logfile)
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
	return log.New(LogWriter(),
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
