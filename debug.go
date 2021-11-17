package augo

import (
	"fmt"
	"path/filepath"
	"strings"
)

var debugTitle = fmt.Sprintf("[%s-Debug]", LogTitle)

func setDeBugTitle(title string) {
	debugTitle = fmt.Sprintf("[%s-Debug]", title)
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

var DebugPrintRouteFunc func(absolutePath string, nuHandlers int)

func debugPrintRoute(absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		if DebugPrintRouteFunc == nil {
			debugPrint("[SERVICE] %-6s --> %-4s (%d handlers)\n", absolutePath, filepath.Base(absolutePath), nuHandlers)
		} else {
			DebugPrintRouteFunc(absolutePath, nuHandlers)
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
