package template

import (
	"embed"
	"html/template"
	"io"
)

//go:embed views/*.html
var viewsFs embed.FS

type TemplateEngine struct {
	tmpl *template.Template
}

func NewTemplateEngine() (*TemplateEngine, error) {
	tmpl, err := template.ParseFS(viewsFs, "views/*.html")
	if err != nil {
		return nil, err
	}
	return &TemplateEngine{tmpl: tmpl}, nil
}

func (t *TemplateEngine) Execute(w io.Writer, name string, data interface{}) error {
	tmpl, err := t.tmpl.Clone()
	if err != nil {
		return err
	}
	tmpl, err = tmpl.ParseFS(viewsFs, "views/"+name)
	return tmpl.ExecuteTemplate(w, name, data)
}
