package controllers

import (
	"net/http"
)

type Users struct {
	Template struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Template.New.Execute(w, nil)
}
