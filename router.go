package venonat

import "path"

//import "net/http"

type (
	IRouter interface {
		IRoutes
		//Group(string, ...HandlerFunc) *RouterGroup
	}

	IRoutes interface {
		Use(...HandlerFunc) IRoutes

		//Handle(string, string, ...HandlerFunc) IRoutes
		//Any(string, ...HandlerFunc) IRoutes
		GET(string, ...HandlerFunc) IRoutes
		//POST(string, ...HandlerFunc) IRoutes
		//DELETE(string, ...HandlerFunc) IRoutes
		//PATCH(string, ...HandlerFunc) IRoutes
		//PUT(string, ...HandlerFunc) IRoutes
		//OPTIONS(string, ...HandlerFunc) IRoutes
		//HEAD(string, ...HandlerFunc) IRoutes
		//
		//StaticFile(string, string) IRoutes
		//Static(string, string) IRoutes
		//StaticFS(string, http.FileSystem) IRoutes
	}

	RouterGroup struct {
		Handlers HandlersChain
		basePath string
		engine   *Engine
		root     bool
	}
)

var _ IRouter = &RouterGroup{}

func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle("GET", relativePath, handlers)

	return group.returnObj()
}

func (group *RouterGroup) handle(method string, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(group.basePath, relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(method, absolutePath, handlers)

	return group.returnObj()
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergeHandlers := make(HandlersChain, finalSize)
	copy(mergeHandlers, group.Handlers)
	copy(mergeHandlers[len(group.Handlers):], handlers)
	return mergeHandlers
}

func (group *RouterGroup) calculateAbsolutePath(basePath, relativePath string) string {
	return joinPaths(basePath, relativePath)
}

func joinPaths(absolutePath, relativePath string) string {
	if len(relativePath) == 0 {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
