package framework

import (
	"context"
	"net/http"
	"sync"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter

	//写保护机制
	writerMux *sync.Mutex

	// 超时标记
	hasTimeout bool

	//当前请求的handler链条
	handlers []ControllerHandler
	index    int //当前请求调用到调用链的哪个节点,初始值为-1

	// 路由参数map
	params map[string]string

	// Container
	LadeContainer
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

//对外暴露锁
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func NewContext(request *http.Request, response http.ResponseWriter) *Context {
	return &Context{
		request:    request,
		response:   response,
		writerMux:  &sync.Mutex{},
		hasTimeout: false,
		index:      -1,
	}
}

// func (ctx *Context) WithTimeout(parent *Context, timeout time.Duration) (*Context, context.CancelFunc) {
// 	ctxTimeout, cancelFunc := context.WithTimeout(ctx.BaseContext(), timeout)
// 	ctx.request = ctx.request.WithContext(ctxTimeout)

// 	return ctx, cancelFunc
// }

// 调用context的handler链的下一个函数
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}

// 设置handlers
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

// 设置路由参数map，params
func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.response
}
