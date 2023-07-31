package controllers

import (
	"fmt"
	"net/http"
	"web_dev/models"
)

type Users struct {
	Template struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.SignIn.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: For a long term we should show a warning about not being able to sign in
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	// fmt.Fprintf(w, "User created: %+v", user)
	// fmt.Fpr int(w, "Password: ", r.FormValue("password"))
}

func (u Users) ProccessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie(CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}
