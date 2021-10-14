package augo

//為系統在的掃描路徑 給個路徑都會有自己的Handlers
type Nodes map[string]HandlersChain

func (n Nodes) IsExist(root string) bool {
	_, e := n[root]
	return e
}

func (n Nodes) IsEmpty() bool {
	return len(n) > 0
}

//設置該Node的路徑以及對應的Handlers
func (n Nodes) Set(root string, handlers HandlersChain) {
	n[root] = append(n[root], handlers...)
}

//取得Node的Handlers
func (n Nodes) Get(root string) HandlersChain {
	return n[root]
}
