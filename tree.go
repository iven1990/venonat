package venonat

type (
	methodTree struct {
		method string
		nodes  nodes
	}

	methodTrees []*methodTree

	node struct {
		handlers HandlersChain
		path     string
	}

	nodes []*node
)

func (trees methodTrees) get(method string) *methodTree {
	for _, tree := range trees {
		if tree.method == method {
			return tree
		}
	}
	return nil
}

//addRoute 添加路由和对应的方法链
func (node *node) addRoute(path string, handlers HandlersChain) {
	node.path = path
	node.handlers = handlers
}

//getValue 获取路由对应的方法链
func (nodes nodes) getValue(path string) HandlersChain {
	for _, node := range nodes {
		if path == node.path {
			return node.handlers
		}
	}
	return nil
}

func (nodes nodes) addNode (node *node) {
	nodes = append(nodes, node)
}