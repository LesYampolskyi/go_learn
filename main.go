package main

import (
	"log"
	"net/http"
	"path/filepath"
	"web_dev/views"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tpl := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tpl, nil)

}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tpl := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tpl, nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tpl := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tpl, nil)
}

func executeTemplate(w http.ResponseWriter, filepath string, data interface{}) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was template parsing error", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

type User struct {
	Name string
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user := User{Name: username}
	tpl := filepath.Join("templates", "user.gohtml")
	log.Println(user)
	executeTemplate(w, tpl, user)

}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/user/{username}", userHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	http.ListenAndServe(":8080", r)
}
