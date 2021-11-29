package augo

import (
	"os"
	"path/filepath"
	"strings"
)

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

	if lastChar(relativePath) == versionchar && lastChar(finalpath) != versionchar {
		return finalpath + pathChar
	}
	return finalpath
}

func deletFiles(path []string) error {
	for _, p := range path {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

func getmethod(dir string) string {
	return dir[strings.LastIndex(dir, pathChar)+1:]
}
