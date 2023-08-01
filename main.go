package main

import (
	"fmt"
	"net/http"
	"web_dev/controllers"
	"web_dev/models"
	"web_dev/templates"
	"web_dev/views"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

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

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INT UNIQUE,
		token_hash TEXT NOT NULL
	);
`)

	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
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
	r.Post("/sign-in", usersC.ProcessSignIn)
	r.Post("/sign-out", usersC.ProcessingSignOut)
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
