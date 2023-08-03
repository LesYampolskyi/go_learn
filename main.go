package main

import (
	"fmt"
	"net/http"
	"web_dev/controllers"
	"web_dev/migrations"
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
	// Setup database
	cfg := models.GetDefaultPostgresConfig()

	db, err := models.Open(cfg)
	fmt.Println(cfg.String())

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, "")
	if err != nil {
		panic(err)
	}

	// Set up services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	crfKey := "Lb38utheFHaXKAMD6pOgkAHyrPeA2nZV"
	crfMw := csrf.Protect([]byte(crfKey),
		csrf.Secure(false))

	// Set up controllers
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

	// Setup router and routes
	r := chi.NewRouter()
	r.Use(crfMw)
	r.Use(umw.SetUser)
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
	r.Get("/sign-up", usersC.New)
	r.Get("/sign-in", usersC.SignIn)
	r.Post("/sign-in", usersC.ProcessSignIn)
	r.Post("/sign-out", usersC.ProcessingSignOut)
	r.Post("/users", usersC.Create)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello")
		})
	})
	r.Get("/users/me", usersC.CurrentUser)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting the server on port :8080...")
	http.ListenAndServe(":8080", r)
}
