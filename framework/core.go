package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router      map[string]*Tree //二级map匹配路由，第一级string匹配方法，第二级string匹配uri
	middlewares []ControllerHandler
}

func NewCore() *Core {
	// 各方法对应一个Trie Tree

	router := map[string]*Tree{}
	router["GET"] = newTree()
	router["POST"] = newTree()
	router["PUT"] = newTree()
	router["DELETE"] = newTree()

	return &Core{router: router, middlewares: make([]ControllerHandler, 0)}
}

func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 封装自定义context
	ctx := NewContext(r, w)

	// 寻找路由
	handlers := c.FindRouteByRequest(r)
	if handlers == nil {
		ctx.SetStatus(404).Json("not found")
		return
	}
	node := c.FindRouteNodeByRequest(r)
	// 设置路由参数
	params := node.parseParamsFromEndNode(r.URL.Path)
	ctx.SetParams(params)

	// 设置context中的handlers
	ctx.SetHandlers(handlers)

	// 调用路由函数，如果返回err，代表存在内部错误，返回500代码
	if err := ctx.Next(); err != nil {
		ctx.SetStatus(500).Json("inner error")
		return
	}

}

// 注册各个方法的路由
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	if err := c.router["GET"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	if err := c.router["POST"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	if err := c.router["PUT"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(url, handlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	// 查找第一层
	if methodHandlers, ok := c.router[upperMethod]; ok {
		handlers := methodHandlers.FindHandler(uri)
		return handlers
	}

	return nil

}

// 设置中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = middlewares
}

// 寻找匹配最终节点
func (c *Core) FindRouteNodeByRequest(request *http.Request) *node {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)
	if tree, ok := c.router[upperMethod]; ok {
		n := tree.root.matchNode(uri)
		return n
	}
	return nil
}
