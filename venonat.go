package venonat

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

type H map[string]interface{}
type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

type (
	RouteInfo struct {
		Method  string
		Path    string
		Handler string
	}

	Engine struct {
		RouterGroup
		pool  sync.Pool
		trees methodTrees
	}
)

func New() *Engine {
	e := &Engine{
		RouterGroup: RouterGroup{
			basePath: "/",
			root:     true,
			Handlers: nil,
		},
		trees: make(methodTrees, 0, 9),
	}

	e.RouterGroup.engine = e
	e.pool.New = func() interface{} {
		return e.allocateContext()
	}
	return e
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.responseWriter.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(context *Context) {
	httpMethod := context.Request.Method
	path := context.Request.URL.Path

	trees := engine.trees
	for i, length := 0, len(engine.trees); i < length; i++ {
		if trees[i].method == httpMethod {
			fmt.Println(path)
			tree := trees.get(httpMethod)
			handlers := tree.nodes.getValue(path)

			if handlers != nil {
				context.handlers = handlers
				context.Next()
				context.responseWriter.WriteHeaderNow()
				return
			}
		}
	}

}

//将路由添加到树结构中
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	fmt.Println(method, path)

	tree := engine.trees.get(method)

	if tree == nil {
		nodes := make([]*node, 0)
		tree = &methodTree{method: method, nodes: nodes}
		engine.trees = append(engine.trees, tree)
	}

	node := new(node)
	node.addRoute(path, handlers)

	tree.nodes = append(tree.nodes, node)
}

func (engine *Engine) Run(addr ...string) error {

	address := resolveAddress(addr)
	return http.ListenAndServe(address, engine)
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}
