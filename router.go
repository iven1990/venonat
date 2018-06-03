package venonat

import (
	"os"
	"path/filepath"
	"strings"
)

type (
	IRoutes interface {
		Use(...HandlerFunc) IRoutes

		GET(string, ...HandlerFunc) IRoutes
	}

	RouterGroup struct {
		Handlers HandlersChain
		basePath string
		engine   *Engine
	}
)

func NewGroup(path string, engine *Engine) *RouterGroup {
	if path == "" {
		path = "/"
	}
	return &RouterGroup{
		Handlers: make(HandlersChain, 0),
		basePath: path,
		engine:   engine,
	}
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle("GET", relativePath, handlers)

	return group
}

func (group *RouterGroup) Static(relativePath string, dir string) IRoutes {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if relativePath == "/" {
			relativePath = ""
		}

		trimPath := strings.TrimLeft(dir, "./")
		routePath := relativePath + "/" + strings.TrimPrefix(path, trimPath+"/")
		cHandler := func(newPath string) HandlerFunc {
			return func(c *Context) {
				c.File(newPath)
			}
		}(path)

		group.GET(routePath, cHandler)
		return nil
	})
	return group
}

func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle("POST", relativePath, handlers)
	return group
}

func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle("PATCH", relativePath, handlers)
	return group
}

func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle("DELETE", relativePath, handlers)
	return group
}

func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle("PUT", relativePath, handlers)
	return group
}

func (group *RouterGroup) handle(method string, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.basePath
	if absolutePath == "/" {
		absolutePath = relativePath
	} else {
		absolutePath = absolutePath + relativePath
	}
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(method, absolutePath, handlers)

	return group
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergeHandlers := make(HandlersChain, finalSize)
	copy(mergeHandlers, group.Handlers)
	copy(mergeHandlers[len(group.Handlers):], handlers)
	return mergeHandlers
}
