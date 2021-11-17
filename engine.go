package augo

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type EngineOptions func(*Engine)

//設置最大線程數
func MaxThread(threadcount int) EngineOptions {
	return func(e *Engine) {
		e.maxThread = threadcount
	}
}

//掃描資料夾的間隔時間
func ScanIntval(intval time.Duration) EngineOptions {
	return func(e *Engine) {
		e.intval = intval
	}
}

//設置中斷Context
func SetContext(ctx context.Context) EngineOptions {
	return func(e *Engine) {
		e.ctx = ctx
	}
}

//設置Collector
func SetCollector(c *Collector) EngineOptions {
	return func(e *Engine) {
		e.C = c
	}
}

type Engine struct {
	maxThread int
	intval    time.Duration
	s         *Scheduler
	C         *Collector

	scanpath cmap.ConcurrentMap //掃描地址 線程安全
	wg       *sync.WaitGroup
	ctx      context.Context
}

func NewEngine(options ...EngineOptions) *Engine {
	e := &Engine{}
	e.defualtParms()
	for _, opt := range options {
		opt(e)
	}
	e.s.ConfigQueue(e.maxThread)
	return e
}

func (e *Engine) defualtParms() {
	e.maxThread = 3
	e.ctx = context.Background()
	e.s = NewScheduler()
	e.wg = &sync.WaitGroup{}
	e.intval = time.Millisecond * 200
	e.scanpath = cmap.New()
	e.C = NewCollector()
}

func (e *Engine) Wait() {
	e.wg.Wait()
}

func (e *Engine) Run() {
	//設置掃描地址
	e.setScanPaths()
	//背景提交任務
	e.wg.Add(1)
	go e.submit()
	//完成信號
	var complete chan struct{} = make(chan struct{})
	//從適配器中獲取請求並執行
	for i := 0; i < e.maxThread; i++ {
		e.wg.Add(1)
		go e.scheduler(e.s.RequestChan(), complete)
	}
	//開啟適配器
	go e.s.RunByContext(e.ctx, complete)

	debugPrint("Services are driven by %s", GetSystemVersion())
	debugPrint("%d threads used in the background", e.maxThread+2)
}

func (e *Engine) Close() error {
	e.ctx.Done()
	return e.C.Close()
}

//從適配器中獲取請求並執行
func (e *Engine) scheduler(in <-chan *Request, complete chan struct{}) {
	defer e.wg.Done()
	for req := range in {

		if err := e.C.Request(req); err != nil {
			e.C.HandleOnErr(req, err)
		}

		//將所有檔案紀錄 方便記錄去重
		for _, f := range req.Files {
			e.C.Visited(req.method, filepath.Base(f))
		}

		//將完成的請求地址
		e.scanpath.Set(req.root, true)
		//完成信號
		complete <- struct{}{}
	}
}

func (e *Engine) setScanPaths() {
	//確認節點是否為空
	errormessage(e.C.nodes.IsEmpty(), "ScanDir is empty")
	for root := range e.C.nodes {
		e.scanpath.Set(root, true)
	}
}

//定期掃描提交並提交請求
func (e *Engine) submit() {
	defer e.wg.Done()
	t := time.NewTicker(e.intval)
	for {
		select {
		case <-t.C:
			e.s.Submits(e.scanDir()...)
		case <-e.ctx.Done():
			return
		}
	}
}

//掃描資料夾獲取請求
func (e *Engine) scanDir() []*Request {

	var reqs []*Request
	for _, root := range e.scanpath.Keys() {
		b, _ := e.scanpath.Get(root)
		//當該root在處理請求時 跳過掃描
		if !b.(bool) {
			continue
		}

		files, err := ioutil.ReadDir(root)
		if len(files) <= 0 && err == nil {
			continue
		}

		req := NewRequest(root)

		if err != nil {
			e.C.HandleOnErr(req, err)
			continue
		}

		if files, err = e.repeatScan(len(files), root); err != nil {
			e.C.HandleOnErr(req, err)
			continue
		}

		for _, file := range files {

			if e.C.IsVisited(req.method, file.Name()) {
				continue
			}
			req.Files = append(req.Files, filepath.Join(root, file.Name()))
		}

		//檔案均為使用過且尚未刪除檔案
		if len(req.Files) == 0 {
			continue
		}

		reqs = append(reqs, req)
		//當該root有請求存在時 在請求完成時會設置為false
		e.scanpath.Set(root, false)

	}
	return reqs
}

func (e *Engine) repeatScan(filecount int, path string) ([]os.FileInfo, error) {
	//如果查詢到新的檔案 等待 並查詢到所有檔案為止
	time.Sleep(time.Millisecond * 300)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == filecount || len(files) == 0 {
		return files, nil
	}

	return e.repeatScan(len(files), path)
}
