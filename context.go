package augo

import (
	"math"
	"sync"
	"time"
)

//Handlers的最大值
const abortIndex int8 = math.MaxInt8 / 6

//為中間件可以自行定義
type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type ErrMsgs []error

func (msgs ErrMsgs) IsEmpty() bool {
	for _, e := range msgs {
		if e != nil {
			return false
		}
	}
	return true
}

func (msgs ErrMsgs) Error() string {
	s := ""
	for _, msg := range msgs {
		if s != "" {
			s += " ,"
		}
		s += msg.Error()
	}
	return s
}

type Context struct {
	handlers HandlersChain
	Request  *Request
	Keys     map[string]interface{}
	Errs     ErrMsgs
	index    int8
	mu       sync.RWMutex
}

func (c *Context) reset(req *Request) {
	c.Request = req
	c.handlers = nil
	c.Keys = nil
	c.index = -1
	c.Errs = c.Errs[0:0]
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

//中止接下的handlers
func (c *Context) Abort() {
	c.index = abortIndex
}

//中止接下的handlers 並保存err至Errs中
func (c *Context) AbortWithError(err error) error {
	c.Abort()
	return c.Error(err)
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

//保存error至Errs中
func (c *Context) Error(err error) error {
	if err == nil {
		return err
	}
	c.Errs = append(c.Errs, err)
	return err
}

//******************************************************************
//**************************   MataData   **************************
//******************************************************************

//輸入key獲取對應val值
func (c *Context) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exist := c.Keys[key]
	return val, exist
}

//設置kv
func (c *Context) Set(key string, val interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = val
	c.mu.Unlock()
}

//輸入key獲取對應String類型val
func (c *Context) GetString(key string) (s string, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		s, b = val.(string)
	}
	return
}

//輸入key獲取對應Bool類型val
func (c *Context) GetBool(key string) (i bool, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		i, b = val.(bool)
	}
	return
}

//輸入key獲取對應Int類型val
func (c *Context) GetInt(key string) (i int, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		i, b = val.(int)
	}
	return
}

//輸入key獲取對應Int64類型val
func (c *Context) GetInt64(key string) (i64 int64, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, b = val.(int64)
	}
	return
}

//輸入key獲取對應float64類型val
func (c *Context) GetFloat64(key string) (f64 float64, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, b = val.(float64)
	}
	return
}

//輸入key獲取對應Time類型val
func (c *Context) GetTime(key string) (t time.Time, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		t, b = val.(time.Time)
	}
	return
}

//輸入key獲取對應Duration類型val
func (c *Context) GetDuration(key string) (d time.Duration, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		d, b = val.(time.Duration)
	}
	return
}

//輸入key獲取對應StringSlice類型val
func (c *Context) GetStringSlice(key string) (ss []string, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, b = val.([]string)
	}
	return
}

//輸入key獲取對應map[string]interface{}類型val
func (c *Context) GetStringMap(key string) (sm map[string]interface{}, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, b = val.(map[string]interface{})
	}
	return
}

//輸入key獲取對應map[string]string類型val
func (c *Context) GetStringMapString(key string) (sms map[string]string, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, b = val.(map[string]string)
	}
	return
}

//輸入key獲取對應map[string][]string類型val
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string, b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, b = val.(map[string][]string)
	}
	return
}
