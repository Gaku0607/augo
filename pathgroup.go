package augo

import "fmt"

type IPaths interface {
	Use(...HandlerFunc) IPaths
	Handler(string, ...HandlerFunc) IPaths
}

type PathGroup struct {
	basepath  string
	root      bool
	handlers  HandlersChain
	collector *Collector
}

//添加該Group的Handlers
func (g *PathGroup) Use(middlefunc ...HandlerFunc) IPaths {
	g.handlers = append(g.handlers, middlefunc...)
	return g.retrunObj()
}

//從該Group中衍生出新的Group 將會繼承源Group的路徑以及Handlers
func (g *PathGroup) Group(basepath string, handlers ...HandlerFunc) *PathGroup {
	return &PathGroup{
		handlers:  g.combineHandlers(handlers),
		basepath:  g.calculateAbsolutePath(basepath),
		collector: g.collector,
	}
}

//將Group註冊到Collector中
func (g *PathGroup) Handler(method string, handlers ...HandlerFunc) IPaths {
	g.collector.addPaths(
		method,
		g.calculateAbsolutePath(method),
		g.combineHandlers(handlers),
	)
	return g.retrunObj()
}

//返回該Group的絕對路徑
func (g *PathGroup) BasePath() string {
	return g.basepath
}

//合併Handlers
func (g *PathGroup) combineHandlers(handlers []HandlerFunc) HandlersChain {
	handlerSize := len(g.handlers) + len(handlers)
	if handlerSize > int(abortIndex) {
		panic(fmt.Sprintf("too many handlers With %s", g.basepath))
	}

	mergehandler := make(HandlersChain, handlerSize)
	copy(mergehandler, g.handlers)
	copy(mergehandler[len(g.handlers):], handlers)
	return mergehandler
}

//合併AbsolutPath
func (g *PathGroup) calculateAbsolutePath(relativepath string) string {
	return joinPaths(g.basepath, relativepath)
}

func (g *PathGroup) retrunObj() IPaths {
	if g.root {
		return g.collector
	}
	return g
}
