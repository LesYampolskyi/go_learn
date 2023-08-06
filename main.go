package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"web_dev/controllers"
	"web_dev/migrations"
	"web_dev/models"
	"web_dev/templates"
	"web_dev/views"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type User struct {
	Name string
}

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.GetDefaultPostgresConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}

	cfg.CSRF.Key = "Lb38utheFHaXKAMD6pOgkAHyrPeA2nZV"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":3000"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup database

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, "")
	if err != nil {
		panic(err)
	}

	// Set up services
	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	// Set up middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	crfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

	// Set up controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	usersC.Template.New = views.Must(views.ParseFS(
		templates.FS,
		"sign-up.gohtml", "tailwind.gohtml",
	))
	usersC.Template.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"sign-in.gohtml", "tailwind.gohtml",
	))
	usersC.Template.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
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
	fmt.Printf("Starting the server on port :%s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
