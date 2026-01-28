package hb

/*

  configurable CLI debugging

  debug := hb.Debug("prefix", my_debug_level)         << output to STDERR
  debug := hb.DebugLog("prefix", my_debug_level)      << output to log with timestamp

  debug(level, "message")
  debug(level, "Sprintf format string", ...arguments)

  level:
    -1 - Show(unconditionally) + panic !!   == ALERT/Emergency
     0 - Show(unconditionally)              == Important Error
     1 - Show if debug_level >= 1
     n - Show if debug_level >= n

  level suggestion:
   -1 : CRIT (FATAL)                                                 ~= LOG_CRIT
    0 : ERR                                                          ~= LOG_ERROR
    1 : important debug, some modules may always have this debug on  ~= LOG_WARNING level
    2 : less important messages, maybe repeating every minute        ~= LOG_NOTICE
    3 : debug stuff                                                  ~= LOG_INFO
    4 : debug stuff. may give lots of messages at once               ~= LOG_DEBUG


 // TODO - add colors for STDERR output:
 // TODO DebugColor
 // TODO DebugCli

*/

import (
	"fmt"
	"log"
	"os"
)

// DebugFunction is a function type for debug output
type DebugFunction func(level int, format string, a ...any)

// Debug creates a debug function that outputs to STDERR
func Debug(prefix string, level int) DebugFunction {
	return func(l int, format string, a ...any) {
		if l > level {
			return
		}
		p := ""
		switch l {
		case 0:
			p = "ERROR "
		case -1:
			p = "FATAL "
		}
		if a == nil {
			fmt.Fprintf(os.Stderr, "%s%s %s\n", p, prefix, format)
		} else {
			fmt.Fprintf(os.Stderr, "%s%s %s\n", p, prefix, fmt.Sprintf(format, a...))
		}
		if l == -1 {
			log.Panic(prefix + " FATAL")
		}
	}
}

// DebugLog creates a debug function with timestamp prefix (uses log package)
func DebugLog(prefix string, level int) DebugFunction {
	return func(l int, format string, a ...any) {
		if l > level {
			return
		}
		p := ""
		switch l {
		case 0:
			p = "ERROR "
		case -1:
			p = "FATAL "
		}
		if a == nil {
			log.Printf("%s%s %s\n", p, prefix, format)
		} else {
			log.Printf("%s%s %s\n", p, prefix, fmt.Sprintf(format, a...))
		}
		if l == -1 {
			log.Panic(prefix + " FATAL")
		}
	}
}
