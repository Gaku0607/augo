package augo

import "sync"

type VisitStorage interface {
	IsVisited(string, string) bool
	Visited(string, string)
}

type HasStorage struct {
	rw         sync.RWMutex
	visitedmap map[uint64]bool
}

func defaultVisitStorage() VisitStorage {
	return &HasStorage{rw: sync.RWMutex{}, visitedmap: make(map[uint64]bool)}
}

func (h *HasStorage) IsVisited(method, filename string) bool {
	h.rw.RLock()
	defer h.rw.RUnlock()
	if ok := h.visitedmap[hasCode(method, filename)]; ok {
		return true
	}
	return false
}

func (h *HasStorage) Visited(method, filename string) {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.visitedmap[hasCode(method, filename)] = true
}
