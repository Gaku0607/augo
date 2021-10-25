package augo

import (
	"fmt"
	"path/filepath"
	"strings"
)

var debugTitle = fmt.Sprintf("[%s-Debug]", LogTitle)

func IsDebugging() bool {
	return debugCode == augomode
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
