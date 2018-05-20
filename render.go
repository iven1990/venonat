package venonat

import "html/template"

func (e *Engine) LoadHtmlGlob(path string) {
	templ := template.Must(template.New("").Delims("{{", "}}").Funcs(template.FuncMap{}).ParseGlob(path))
	e.htmlTmpl = templ
}

func (c *Context) HTML(path string, data interface{}) {
	c.engine.htmlTmpl.ExecuteTemplate(c.Writer, path, data)
}