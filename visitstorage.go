package augo

import "sync"

type VisitStorage interface {
	IsVisited(string) bool
	Visited(string)
	Size() int
	GetVisited() (string, bool)
	RemoveVisited(string)
}

type PathStorage struct {
	rw         sync.RWMutex
	visitedmap map[string]bool
	size       int
}

func defaultVisitStorage() VisitStorage {
	return &PathStorage{rw: sync.RWMutex{}, visitedmap: make(map[string]bool)}
}

func (h *PathStorage) IsVisited(path string) bool {
	h.rw.RLock()
	defer h.rw.RUnlock()
	if ok := h.visitedmap[path]; ok {
		return true
	}
	return false
}

func (h *PathStorage) Visited(path string) {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.visitedmap[path] = true
	h.size++
}

func (h *PathStorage) RemoveVisited(path string) {
	h.rw.Lock()
	defer h.rw.Unlock()
	if exist := h.visitedmap[path]; !exist {
		return
	}
	delete(h.visitedmap, path)
	h.size--
	return
}

func (h *PathStorage) GetVisited() (path string, exist bool) {
	h.rw.RLock()
	defer h.rw.RUnlock()
	if h.size <= 0 {
		return
	}

	for key := range h.visitedmap {
		return key, true
	}
	return
}

func (h *PathStorage) Size() int {
	h.rw.RLock()
	defer h.rw.RUnlock()
	return h.size
}
