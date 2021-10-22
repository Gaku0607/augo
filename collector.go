package augo

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Gaku0607/augo/logger"
)

type CollectorOption func(*Collector)

//修改CollectorGroup中的根路徑
func CollectorGroupRoot(root string) CollectorOption {
	return func(c *Collector) {
		c.PathGroup.basepath = root
	}
}

//修改Collector的LogMod
func LoggerMode(b bool) CollectorOption {
	return func(c *Collector) {
		c.LogMode = b
	}
}

//修改Collector的Logger
func LoggerConfig(l *logger.Logger) CollectorOption {
	return func(c *Collector) {
		c.Logger = l
	}
}

//修改Collector的RequestLogKey
func RequestLogKey(f func(*Request) logger.LogKey) CollectorOption {
	return func(c *Collector) {
		c.requestlogkey = f
	}
}

//修改Collector的ErrorLogKey
func ErrorLogKey(f func(*Request, error) logger.LogKey) CollectorOption {
	return func(c *Collector) {
		c.errlogkey = f
	}
}

//修改Collector的ResultLogKey
func ResultLogKey(f func(*Context) logger.LogKey) CollectorOption {
	return func(c *Collector) {
		c.resultlogkey = f
	}
}

var (
	//Collector默認所使用的ErrLogKey
	defaultErrLogKey = func(_ *Request, e error) logger.LogKey {
		return logger.LogKey{"errMsg": e.Error()}
	}
	//Collector默認所使用的ReqLogKey
	defaultReqLogKey = func(r *Request) logger.LogKey {
		return logger.LogKey{"FileCount": len(r.Files)}
	}
	//Collector默認所使用的ResultLogKey
	defaultResultLogKey = func(c *Context) logger.LogKey {
		return c.Keys
	}
)

type Collector struct {
	PathGroup

	//LoggerMode為true時會輸出Log
	//Collector默認開啟Logger

	LogMode bool

	//Collector默認使用 logger.DefaultLogger()
	//調用CollectorOption LoggerConfig 可自行修改輸出格式
	//(參考logger.Logger logger.LoggerConfig)

	Logger *logger.Logger

	//Log輸出REQUEST時 可自定義RequestLogKey的部分
	//Collector默認使用defaultRequestLogKey

	requestlogkey func(*Request) logger.LogKey

	//Log輸出ERROR時 可自定義ErrLogKey的部分
	//Collector默認使用defaultErrLogKey

	errlogkey func(*Request, error) logger.LogKey

	//Log輸出Result時 可自定義ResultLogKey的部分
	//Collector默認使用defaultResultLogKey

	resultlogkey func(*Context) logger.LogKey

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
func DefautCollector() *Collector {
	c := &Collector{}
	c.defautParms()
	c.Use(Recovery(c.Logger), DeletFiles())
	return c
}

func (c *Collector) defautParms() {
	c.PathGroup = PathGroup{
		basepath:  "",
		root:      true,
		collector: c,
	}
	c.nodes = make(map[string]HandlersChain)
	c.pool.New = func() interface{} {
		return c.allocateContext()
	}
	c.LogMode = true
	c.Logger = logger.NowLogger()
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
	ctx.handlers = c.nodes[ctx.Request.root]
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
func (c *Collector) addPaths(basepath string, handlers []HandlerFunc) {
	errormessage(len(basepath) > 0, "Path can not be empty")
	errormessage(len(handlers) > 0, "Handlers can not be empty")

	errormessage(
		!c.nodes.IsExist(basepath),
		fmt.Sprintf("%s is exist", basepath),
	)
	c.nodes.Set(basepath, handlers)
}

//當LoggerMode為true時會調用指定的Logger以及Logkey
func (c *Collector) HandleOnErr(r *Request, err error) {
	if c.LogMode {
		Parms := logger.CreateLogParms(r.Id, logger.ERROR, r.FilesName(), r.Method(), c.errlogkey(r, err))
		c.Logger.Log(Parms)
	}
}

//當LoggerMode為true時會調用指定的Logger以及Logkey
func (c *Collector) handleOnRequest(r *Request) {
	if c.LogMode {
		Parms := logger.CreateLogParms(r.Id, logger.REQUEST, r.FilesName(), r.Method(), c.requestlogkey(r))
		c.Logger.Log(Parms)
	}
}

//當LoggerMode為true時會調用指定的Logger以及Logkey
func (c *Collector) handleOnResult(ctx *Context) {
	if c.LogMode {
		Parms := logger.CreateLogParms(ctx.Request.Id, logger.RESULT, ctx.Request.FilesName(), ctx.Request.Method(), c.resultlogkey(ctx))
		c.Logger.Log(Parms)
	}
}

func (c *Collector) allocateContext() interface{} {
	return &Context{}
}
