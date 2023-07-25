package main

import (
	"fmt"
	"log"
	"net/http"
	"web_dev/controllers"
	"web_dev/templates"
	"web_dev/views"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

// func userHandler(w http.ResponseWriter, r *http.Request) {
// 	username := chi.URLParam(r, "username")
// 	user := User{Name: username}
// 	tpl := filepath.Join("templates", "user.gohtml")
// 	log.Println(user)
// 	executeTemplate(w, tpl, user)
// }

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml"))))

	// r.Get("/user/{username}", userHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on port :8080....")
	http.ListenAndServe(":8080", r)
}
