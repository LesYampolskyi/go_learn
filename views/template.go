package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"web_dev/context"
	"web_dev/models"

	"github.com/gorilla/csrf"
)

type Template struct {
	htmlTlp *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tpl := template.New(pattern[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField is not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser is not implemented")
			},
			"errors": func() []string {
				return []string{
					"Don't do that",
					// "The email address you provided is already associated with an account",
					"Something went wrong",
				}
			},
		},
	)
	tpl, err := tpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template %v", err)
	}
	return Template{
		htmlTlp: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// TODO:: check without Cloning
	tpl, err := t.htmlTlp.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)
	// var buf bytes.Buffer
	// err = tpl.Execute(&buf, data)
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was template executing error", http.StatusInternalServerError)
		return
	}
	// io.Copy(w, &buf)
}
