package framework

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Use(...ControllerHandler)
}

type Group struct {
	core        *Core
	prefix      string
	middlewares []ControllerHandler
	parent      *Group
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Post(uri, allHandlers...)
}
func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Put(uri, allHandlers...)
}
func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Delete(uri, allHandlers...)
}

// 从core中初始化Group
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (g *Group) Group(prefix string) IGroup {
	subGroup := NewGroup(g.core, g.prefix+prefix)
	subGroup.parent = g
	return subGroup
}

// Group设置中间件
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = middlewares
}

// Group获取中间件
// Group内使用父路由的中间件+自身中间件
func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}

	return append(g.parent.getMiddlewares(), g.middlewares...)
}
