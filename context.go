package venonat

import (
	"net/http"
	"github.com/iven1990/enheng/venonat/render"
)

type Context struct {
	engine         *Engine
	responseWriter responseWriter

	Request *http.Request
	Writer  ResponseWriter

	handlers HandlersChain

	index int8
}

func (c *Context) reset() {
	c.Writer = &c.responseWriter
	c.handlers = nil
	c.index = -1
}

func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++  {
		c.handlers[c.index](c)
	}
}

func (c *Context) Status(code int) {
	c.responseWriter.WriteHeader(code)
}

func (c *Context) JSON (code int, obj interface{}) {
	c.Status(code)
	if err := render.WriteJSON(c.Writer, obj); err != nil {
		panic(err)
	}
}
