package augo

import (
	"context"
	"time"
)

type SchedulerOptions func(*Scheduler)

//設置休眠時間
func SleepIntvar(t time.Duration) SchedulerOptions {
	return func(s *Scheduler) {
		s.sleepIntvar = t
	}
}

//設置提交間格時間
func SubmitIntvar(t time.Duration) SchedulerOptions {
	return func(s *Scheduler) {
		s.submitIntvar = time.NewTicker(t)
	}
}

//設置請求的儲存格式
func SetRequestStore(store RequestStore) SchedulerOptions {
	return func(s *Scheduler) {
		s.RequestStore = store
	}
}

type Scheduler struct {
	RequestStore //儲存請求 默認使用Momory

	sleepIntvar  time.Duration // 當系統進入閒置狀態時 休眠
	submitIntvar *time.Ticker  // 提交至上限時 休眠

	threadQueue []chan *Request //存放所有Thread所對應的 RequestChan
	queuesize   int             //ThreadPool的大小
	ptr         int             //紀錄ThreadPool當前所使用的Thread的位子
}

func NewScheduler(opts ...SchedulerOptions) *Scheduler {
	s := &Scheduler{}
	s.defaultParms()
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Scheduler) defaultParms() {
	s.sleepIntvar = time.Millisecond * 1500
	s.submitIntvar = time.NewTicker(time.Millisecond * 1500)
	s.RequestStore = NewRequestStore()
}

//初始化每個線程所使用的ＣＨＡＮ隊列
func (s *Scheduler) ConfigQueue(size int) {
	s.queuesize = size
	s.threadQueue = make([]chan *Request, size)
	for i := 0; i < s.queuesize; i++ {
		s.threadQueue[i] = make(chan *Request)
	}
}

//返回每個線程所對應的請求ＣＨＡＮ
func (s *Scheduler) RequestChan() chan *Request {
	return s.pull()
}

//儲存多個請求
func (s *Scheduler) Submits(reqs ...*Request) {
	for _, req := range reqs {
		s.Submit(req)
	}
}

//儲存請求
func (s *Scheduler) Submit(req *Request) {
	for !s.PushRequest(req) {
		<-s.submitIntvar.C
	}
}

func (s *Scheduler) RunByContext(ctx context.Context, complete <-chan struct{}) {

	var (
		active int //在運行的線程數
		req    *Request
	)
	ticker := time.NewTicker(s.sleepIntvar)
	for {
		var activeThread chan *Request
		//當已經沒有請求 並且 無活動的thread時 代表閒置狀態 進行休眠
		if s.isRequsetEmpty() && active == 0 {
			<-ticker.C //休眠
		}

		if !s.isRequsetEmpty() {
			activeThread = s.pull()
			req = s.PullRequest()
		}

	Loop:
		for {
			select {
			case activeThread <- req:
				active++
				break Loop
			case <-complete:
				active--
				if activeThread == nil && active == 0 {
					break Loop
				}
			case <-ctx.Done():
				//等待所有線程工作完畢才進行ＤＯＮＥ
				for i := active; i > 0; i-- {
					<-complete
				}
				s.closeThread()
				return
			default:
				break Loop
			}
		}
	}
}

func (s *Scheduler) isRequsetEmpty() bool {
	return s.RequestStore.Size() == 0
}

//返回每個線程所對應chan 並 移動ptr
func (s *Scheduler) pull() chan *Request {
	c := s.threadQueue[s.ptr]
	s.ptr = (s.ptr + 1) % s.queuesize
	return c
}

func (s *Scheduler) closeThread() {
	for _, c := range s.threadQueue {
		close(c)
	}
}
