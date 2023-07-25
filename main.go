package main

import (
	"fmt"
	"log"
	"net/http"
	"web_dev/controllers"
	"web_dev/templates"
	"web_dev/views"

	"github.com/go-chi/chi/v5"
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

func main() {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)

	fileServer := http.FileServer(http.Dir("./css/"))

	r.Handle("/css/*", http.StripPrefix("/css/", fileServer))

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml", "tailwind.gohtml",
	))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"contact.gohtml", "tailwind.gohtml",
	))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml", "tailwind.gohtml",
	))))

	usersC := controllers.Users{}

	usersC.Template.New = views.Must(views.ParseFS(
		templates.FS,
		"sign-up.gohtml", "tailwind.gohtml",
	))

	r.Get("/sign-up", usersC.New)

	// r.Get("/user/{username}", userHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on port :8080...")
	http.ListenAndServe(":8080", r)
}
