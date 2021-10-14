package augo

import (
	"sync"
)

type RequestStore interface {
	//添加Request至Storage
	PushRequest(*Request) bool
	//從Storage取得Request
	PullRequest() *Request
	//返回當前儲存的個數
	Size() int
}

type RequestNode struct {
	Request *Request
	next    *RequestNode
}

type MemoryQueueOptions func(*InMemoryRequestQueue)

//設置對大儲存數 為0時無上限值
func MaxMemorySize(size int) func(*InMemoryRequestQueue) {
	return func(imrq *InMemoryRequestQueue) {
		imrq.maxSize = size
	}
}

//使用記憶體儲存請求
type InMemoryRequestQueue struct {
	maxSize int          //儲存上限
	size    int          //當前容量
	frist   *RequestNode //Queue的頭指針 Pull時先從frist開始取
	last    *RequestNode //Queue的尾指針 Push從last以後開始添加

	rw *sync.RWMutex
}

func NewRequestStore(opts ...MemoryQueueOptions) *InMemoryRequestQueue {
	s := &InMemoryRequestQueue{}
	s.defaultParms()
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (r *InMemoryRequestQueue) defaultParms() {
	r.maxSize = 1000
	r.rw = &sync.RWMutex{}
}

//傳入Request並將其設為last 當超過上限時進行Panic
func (r *InMemoryRequestQueue) PushRequest(req *Request) bool {
	r.rw.Lock()
	defer r.rw.Unlock()
	if r.size >= r.maxSize && r.maxSize > 0 {
		return false
	}
	node := &RequestNode{Request: req}

	if r.frist == nil {
		r.frist = node
	} else {
		r.last.next = node
	}
	r.last = node
	r.size++
	return true
}

//返回當前儲存容量
func (r *InMemoryRequestQueue) Size() int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.size
}

//推出一個Request 當為0時返回nil
func (r *InMemoryRequestQueue) PullRequest() *Request {
	r.rw.Lock()
	defer r.rw.Unlock()
	if r.size == 0 {
		return nil
	}
	req := r.frist.Request
	r.frist = r.frist.next
	r.size--
	return req
}
