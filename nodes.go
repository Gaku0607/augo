package augo

import "fmt"

//為系統在的掃描路徑 給個路徑都會有自己的Handlers
type Nodes map[string]*Node

func (n Nodes) VisitMode(root string) bool {
	if node, exist := n[root]; exist {
		return node.visitMode
	}
	return false
}

func (n Nodes) IsExist(root string) bool {
	_, e := n[root]
	return e
}

func (n Nodes) IsEmpty() bool {
	return len(n) > 0
}

//設置該Node的路徑以及對應的Handlers
func (n Nodes) Set(root string, handlers HandlersChain, visitmode bool) {
	_, exist := n[root]
	errormessage(!exist, fmt.Sprintf("%s is exist,", root))

	n[root] = &Node{visitMode: visitmode, handlers: handlers}
}

//取得Node的Handlers
func (n Nodes) GetHandlers(root string) HandlersChain {
	return n[root].handlers
}

type Node struct {
	visitMode bool
	handlers  HandlersChain
}
