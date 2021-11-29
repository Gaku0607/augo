package augo

import (
	"fmt"
	"strings"
)

var debugTitle = fmt.Sprintf("[%s]", LogTitle)

func setDeBugTitle(title string) {
	debugTitle = fmt.Sprintf("[%s]", title)
}

func IsDebugging() bool {
	return debugCode == augomode
}

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

var augomode = debugCode
var modeName = DebugMode

// SetMode sets augo mode according to input string.
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		augomode = debugCode
	case ReleaseMode:
		augomode = releaseCode
	case TestMode:
		augomode = testCode
	default:
		panic("augo mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	modeName = value
}

var DebugPrintRouteFunc func(absolutePath string, nuHandlers int, visitmode bool)

func debugPrintRoute(absolutePath string, handlers HandlersChain, visitmode bool) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		if DebugPrintRouteFunc == nil {
			if visitmode {
				debugPrint("[SERVICE][VISITMODE] %-6s -->  (%d handlers)\n", absolutePath, nuHandlers)
			} else {
				debugPrint("[SERVICE] %-6s -->  (%d handlers)\n", absolutePath, nuHandlers)
			}
		} else {
			DebugPrintRouteFunc(absolutePath, nuHandlers, visitmode)
		}
	}
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
	}
	fmt.Fprintf(defaultWriter, debugTitle+format, values...)
}

func debugPrintError(err error) {
	if err != nil {
		if IsDebugging() {
			fmt.Fprintf(defaultErrWriter, debugTitle+"[ERROR] %v\n", err)
		}
	}
}
