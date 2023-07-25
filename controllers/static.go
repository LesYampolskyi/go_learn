package controllers

import (
	"net/http"
	"web_dev/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

//  Nullam nisl tortor, mattis eu consequat a, faucibus sit amet eros. Morbi luctus turpis tellus, at volutpat nisl convallis id. Duis non magna at justo congue elementum in vel nulla.

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Praesent nec lacus odio.",
			Answer:   "Phasellus a urna augue. Nullam eu ex enim. Cras non magna ac eros dapibus rhoncus vel at dui. Donec sit amet suscipit metus.",
		},
		{
			Question: "Phasellus ullamcorper, tortor sit amet varius malesuada",
			Answer:   "Magna elit malesuada ex, nec maximus enim orci quis turpis. Donec sollicitudin ultricies lectus, sed egestas est efficitur at. ",
		},
		{
			Question: "In tempus ante sit amet metus porttitor mattis.",
			Answer:   "Nulla facilisi. Maecenas feugiat tempus ex, vitae sodales est facilisis non. In sed lectus a mi laoreet condimentum. Phasellus vulputate fringilla arcu, nec egestas turpis dignissim ac.",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
