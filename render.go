package venonat

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func (e *Engine) LoadHtmlGlob(path string) {
	templ := template.Must(template.New("").Delims("{{", "}}").Funcs(template.FuncMap{}).ParseGlob(path))
	e.htmlTmpl = templ
}

func (c *Context) HTML(path string, data interface{}) {
	c.engine.htmlTmpl.ExecuteTemplate(c.Writer, path, data)
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) Json(status int, obj interface{}) {
	header := c.Writer.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
	}
	data, _ := json.Marshal(&obj)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)

}
