package main

import (
	"fmt"
	"log"
	"net/http"
	"web_dev/controllers"
	"web_dev/models"
	"web_dev/templates"
	"web_dev/views"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
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

	cfg := models.GetDefaultPostgresConfig()

	db, err := models.Open(cfg)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}

	usersC.Template.New = views.Must(views.ParseFS(
		templates.FS,
		"sign-up.gohtml", "tailwind.gohtml",
	))
	usersC.Template.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"sign-in.gohtml", "tailwind.gohtml",
	))

	r.Get("/sign-up", usersC.New)
	r.Get("/sign-in", usersC.SignIn)
	r.Post("/sign-in", usersC.ProccessSignIn)
	r.Post("/users", usersC.Create)
	r.Get("/users/me", usersC.CurrentUser)

	crfKey := "Lb38utheFHaXKAMD6pOgkAHyrPeA2nZV"
	crfMw := csrf.Protect([]byte(crfKey),
		csrf.Secure(false))
	// r.Get("/user/{username}", userHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on port :8080...")
	http.ListenAndServe(":8080", crfMw(r))
}
