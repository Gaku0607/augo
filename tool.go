package augo

import (
	"hash/fnv"
	"os"
	"path/filepath"
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

	if lastChar(finalpath) != versionchar {
		return finalpath + pathChar
	}
	return finalpath
}

func deletFiles(path []string) error {
	for _, p := range path {
		if err := deleteFile(p); err != nil {
			return err
		}
	}
	return nil
}

func deleteFile(path string) error {
	return os.Remove(path)
}

func hasCode(method, file string) uint64 {
	f := fnv.New64a()
	f.Write([]byte(method))
	f.Write([]byte(file))
	return f.Sum64()
}
