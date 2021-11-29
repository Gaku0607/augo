package augo

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type CollectorOption func(*Collector)

//修改CollectorGroup中的根路徑
func CollectorGroupRoot(root string) CollectorOption {
	return func(c *Collector) {
		c.PathGroup.basepath = root
	}
}

//修改Collector的Logger
func SetLogger(l *Logger) CollectorOption {
	return func(c *Collector) {
		c.Logger = l
	}
}

//修改Collector的RequestLogKey
func RequestLogKey(f func(*Request) LogKey) CollectorOption {
	return func(c *Collector) {
		c.requestlogkey = f
	}
}

//修改Collector的ErrorLogKey
func ErrorLogKey(f func(*Request, error) LogKey) CollectorOption {
	return func(c *Collector) {
		c.errlogkey = f
	}
}

//修改Collector的ResultLogKey
func ResultLogKey(f func(*Context) LogKey) CollectorOption {
	return func(c *Collector) {
		c.resultlogkey = f
	}
}

//修改Collector的VisitStorage
func SetVisitStorage(v VisitStorage) CollectorOption {
	return func(c *Collector) {
		c.visit = v
	}
}

var (
	//Collector默認所使用的ErrLogKey
	defaultErrLogKey = func(_ *Request, e error) LogKey {
		return LogKey{NORMAL_ERROR: e.Error()}
	}
	//Collector默認所使用的ReqLogKey
	defaultReqLogKey = func(r *Request) LogKey {
		return LogKey{"FileCount": len(r.Files)}
	}
	//Collector默認所使用的ResultLogKey
	defaultResultLogKey = func(c *Context) LogKey {
		return c.Keys
	}
)

type Collector struct {
	PathGroup

	//Collector默認使用 logger.DefaultLogger()
	//調用CollectorOption LoggerConfig 可自行修改輸出格式
	//(參考logger.Logger logger.LoggerConfig)

	Logger *Logger

	//

	visit VisitStorage

	//Log輸出REQUEST時 可自定義RequestLogKey的部分
	//Collector默認使用defaultRequestLogKey

	requestlogkey func(*Request) LogKey

	//Log輸出ERROR時 可自定義ErrLogKey的部分
	//Collector默認使用defaultErrLogKey

	errlogkey func(*Request, error) LogKey

	//Log輸出Result時 可自定義ResultLogKey的部分
	//Collector默認使用defaultResultLogKey

	resultlogkey func(*Context) LogKey

	//每個路徑Node對應的Handlers

	nodes Nodes

	//RequestTotal Request的唯一標示

	requestcount int64
	pool         sync.Pool
}

func NewCollector(opts ...CollectorOption) *Collector {
	c := &Collector{}
	c.defautParms()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

//默認使用Recovery and DeletFiles 中間件
func DefautCollector(opts ...CollectorOption) *Collector {
	c := &Collector{}
	c.defautParms()
	for _, opt := range opts {
		opt(c)
	}
	c.Use(Recovery(c.Logger), DeletFiles())

	debugPrint(`[WARNING] Creating an Engine instance with the DeleteFiles and Recovery middleware already attached.

`)
	return c
}

func (c *Collector) defautParms() {
	c.PathGroup = PathGroup{
		basepath:  "",
		root:      true,
		collector: c,
	}
	c.nodes = make(map[string]*Node)
	c.pool.New = func() interface{} {
		return c.allocateContext()
	}

	c.Logger = NowLogger()
	c.visit = defaultVisitStorage()
	c.errlogkey = defaultErrLogKey
	c.requestlogkey = defaultReqLogKey
	c.resultlogkey = defaultResultLogKey
}

//接收請求 找尋對應的Node
func (c *Collector) Request(req *Request) (err error) {
	//設置每個請求的唯一標示
	req.Id = c.setRequestId()

	//打印Log
	c.handleOnRequest(req)

	ctx := c.pool.Get().(*Context)
	//每個req都需重置先前的紀錄
	ctx.reset(req)
	ctx.handlers = c.nodes.GetHandlers(ctx.Request.root)
	ctx.Next()

	if ctx.Errs.IsEmpty() {
		c.handleOnResult(ctx)
	} else {
		err = ctx.Errs
	}

	c.pool.Put(ctx)
	return
}

//每個Request的唯一標示
func (c *Collector) setRequestId() int64 {
	return atomic.AddInt64(&c.requestcount, 1)
}

//添加Handlers至Collector中（根Group）
func (c *Collector) Use(middlefunc ...HandlerFunc) IPaths {
	c.PathGroup.handlers = append(c.PathGroup.handlers, middlefunc...)
	return c
}

func (c *Collector) Close() error {
	return c.Logger.Config.Close()
}

//註冊路徑
func (c *Collector) addPaths(AbsolutePath string, handlers HandlersChain, visitmode bool) {
	errormessage(len(AbsolutePath) > 0, "Path can not be empty")
	errormessage(len(handlers) > 0, "Handlers can not be empty")
	errormessage(!c.nodes.IsExist(AbsolutePath), fmt.Sprintf("%s is exist", AbsolutePath))

	debugPrintRoute(AbsolutePath, handlers, visitmode)

	c.nodes.Set(AbsolutePath, handlers, visitmode)
}

//確認該service下 指定的file是否已被訪問過
func (c *Collector) IsVisited(root, filename string) bool {
	if c.nodes.VisitMode(root) {
		return c.visit.IsVisited(filepath.Join(root, filename))
	}
	return false
}

//將已被訪問過的file儲存 以防重複訪問
func (c *Collector) Visited(root, filename string) {
	if c.nodes.VisitMode(root) {
		c.visit.Visited(filepath.Join(root, filename))
	}
}

//當LoggerMode為true時會調用指定的Logger以及Logkey
func (c *Collector) HandleOnErr(r *Request, err error) {
	Parms := CreateLogParms(r.Id, ERROR, r.FilesName(), r.Method(), c.errlogkey(r, err))
	c.Logger.Log(Parms)
}

//當LoggerMode為true時會調用指定的Logger以及Logkey
func (c *Collector) handleOnRequest(r *Request) {
	Parms := CreateLogParms(r.Id, REQUEST, r.FilesName(), r.Method(), c.requestlogkey(r))
	c.Logger.Log(Parms)
}

//當LoggerMode為true時會調用指定的Logger以及Logkey
func (c *Collector) handleOnResult(ctx *Context) {
	Parms := CreateLogParms(ctx.Request.Id, RESULT, ctx.Request.FilesName(), ctx.Request.Method(), c.resultlogkey(ctx))
	c.Logger.Log(Parms)
}

func (c *Collector) allocateContext() interface{} {
	return &Context{}
}
