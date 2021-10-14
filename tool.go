package augo

import "path/filepath"

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func errormessage(assert bool, text string) {
	if !assert {
		panic(text)
	}
}

func joinPaths(absolutePath, relativePath string) string {
	errormessage(relativePath != "", "relativePath is nil")

	finalpath := filepath.Join(absolutePath, relativePath)
	versionchar := []uint8(pathChar)[0]

	if lastChar(finalpath) != versionchar {
		return finalpath + pathChar
	}
	return finalpath
}
