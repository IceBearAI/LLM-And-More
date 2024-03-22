//go:build !windows && !plan9 && !nacl
// +build !windows,!plan9,!nacl

package logging

import (
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/syslog"
	"log"
	gosyslog "log/syslog"
)

func syslogLogger(lv, serviceName string) kitlog.Logger {
	w, err := gosyslog.New(syslogLevel(lv), serviceName)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return syslog.NewSyslogLogger(w, kitlog.NewLogfmtLogger)
}

func syslogLevel(lv string) gosyslog.Priority {
	switch lv {
	case "emergency":
		return gosyslog.LOG_EMERG
	case "alert":
		return gosyslog.LOG_ALERT
	case "critical":
		return gosyslog.LOG_CRIT
	case "error":
		return gosyslog.LOG_ERR
	case "warning":
		return gosyslog.LOG_WARNING
	case "notice":
		return gosyslog.LOG_NOTICE
	case "info":
		return gosyslog.LOG_INFO
	case "debug":
		return gosyslog.LOG_DEBUG
	default:
		return gosyslog.LOG_LOCAL0
	}
}
