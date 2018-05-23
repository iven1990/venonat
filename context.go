package venonat

import (
	"net/http"
)

type Context struct {
	engine *Engine

	Request *http.Request
	Writer  http.ResponseWriter

	handlers HandlersChain

	index int8
}

func (c *Context) reset() {
	c.handlers = nil
	c.index = -1
}

func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}
