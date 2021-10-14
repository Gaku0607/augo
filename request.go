package augo

import (
	"path/filepath"
	"strings"
)

type Request struct {
	Id     int64
	root   string //資料夾的絕對位址
	method string
	Files  []string
}

func NewRequest(root string, path ...string) *Request {
	dir := filepath.Dir(root)
	return &Request{root: root, Files: path, method: dir[strings.LastIndex(dir, pathChar)+1:]}
}

func (r *Request) FilesName() string {
	var names []string = make([]string, len(r.Files))
	for i, file := range r.Files {
		names[i] = filepath.Base(file)
	}
	return strings.Join(names, " ,")
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Clone() *Request {
	return &Request{Id: r.Id, root: r.root, Files: r.Files}
}
