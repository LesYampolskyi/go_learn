package controllers

import "net/http"

const (
	CookieSession = "session"
)

func newCookie(name, value string) *http.Cookie {
	cookies := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &cookies
}

func setCookie(w http.ResponseWriter, name, value string) {
	cookie := newCookie(name, value)
	http.SetCookie(w, cookie)
}
