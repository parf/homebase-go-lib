package hb

import (
	"log"
	"log/syslog"
	"path/filepath"
	"runtime"
)

// SysLogNotice writes notice to syslog
// Usage: hb.SysLogNotice(fmt.Sprintf(...))
func SysLogNotice(message string) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("Error in syslog, ", message)
	}
	logwriter, _ := syslog.New(syslog.LOG_NOTICE, filepath.Base(filename))
	logwriter.Notice("go-api " + message)
}

// SysLogError writes error to syslog
// Usage: hb.SysLogError(fmt.Sprintf(...))
func SysLogError(message string) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("Error in syslog, ", message)
	}
	logwriter, _ := syslog.New(syslog.LOG_NOTICE, filepath.Base(filename))
	logwriter.Err("go-api " + message)
}
