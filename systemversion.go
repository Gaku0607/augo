package augo

import "fmt"

const (
	Windows = "Windows"
	MacOS   = "MacOS"
	Linux   = "Linux"
)

//默認使用MacOS環境的路徑格式
var system_version = MacOS
var pathChar = `/`
var delete_msg = "no such file"

func SetSystemVersion(version string) {
	switch version {
	case Windows:
		pathChar = `\`
		delete_msg = "The system cannot find the file"

	case MacOS:
		pathChar = `/`
		delete_msg = "no such file"

	case Linux:
		pathChar = `/`
		delete_msg = "no such file"

	default:
		panic(fmt.Sprintf("version:= %s,input version is not format", version))
	}
	system_version = version
}

func GetSystemVersion() string {
	return system_version
}

func GetPathChar() string {
	return pathChar
}
