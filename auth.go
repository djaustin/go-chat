package main

import (
	"log"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if _, err := r.Cookie("Authorization"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		log.Panicln(err)
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func mustAuth(h http.Handler) http.Handler {
	return &authHandler{next: h}
}
