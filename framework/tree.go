package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	isLast  bool
	segment string
	// handler ControllerHandler
	handlers []ControllerHandler //中间件+handler
	childs   []*node

	parent *node
}

// 判断是否是通配段，以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 子节点中查找匹配
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	// 通配模式，则返回所有子节点
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			// 子节点通配模式
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			//静态文本匹配
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// 判断uri是否在节点的子节点存在匹配
func (n *node) matchNode(uri string) *node {
	// 分割成2段
	segments := strings.SplitN(uri, "/", 2)
	// 第一段用于匹配下一层子节点
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配下一层子节点
	cnodes := n.filterChildNodes(segment)
	// 如果没有子节点符合，则直接返回nil
	if len(cnodes) == 0 {
		return nil
	}
	// uri不可再分割，则是最后一个标记
	if len(segments) == 1 {
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		// 不是最后一个节点，但uri已经结束，则不匹配
		return nil
	}

	// 成功分割，说明有子节点查找到，递归查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}

	return nil
}

// 添加路由节点
func (tree *Tree) AddRouter(uri string, handlers ...ControllerHandler) error {
	n := tree.root
	// 确认是否有路由冲突
	if n.matchNode(uri) != nil {
		return errors.New("route exists: " + uri)
	}

	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		islast := index == len(segments)-1

		var objNode *node // 当前匹配的node

		childNodes := n.filterChildNodes(segment)

		// 子节点中找到匹配的node
		if len(childNodes) > 0 {
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			// 创建一个当前node的子节点
			cnode := newNode()
			cnode.segment = segment
			if islast {
				cnode.isLast = true
				cnode.handlers = append(cnode.handlers, handlers...)
			}
			cnode.parent = n
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		// 下一个循环，当前节点指向匹配的节点
		n = objNode

	}
	return nil
}

func newNode() *node {
	return &node{
		segment:  "",
		isLast:   false,
		childs:   make([]*node, 0),
		handlers: make([]ControllerHandler, 0),
	}
}

// 匹配uri，查找对应的handler
func (tree *Tree) FindHandler(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

func newTree() *Tree {
	return &Tree{
		root: newNode(),
	}
}

// 将uri解析为params
func (n *node) parseParamsFromEndNode(uri string) map[string]string {
	ret := map[string]string{}
	segments := strings.Split(uri, "/")
	cnt := len(segments)

	cur := n

	for i := cnt - 1; i >= 0; i-- {
		if cur.segment == "" {
			break
		}
		if isWildSegment(cur.segment) {
			ret[cur.segment[1:]] = segments[i]
		}
		cur = cur.parent
	}
	return ret
}
