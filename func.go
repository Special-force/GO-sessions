package main

import (
	"net/http"

	"github.com/satori/go.uuid"
)

func getus(w http.ResponseWriter, r *http.Request) user {
	c, err := r.Cookie("User-cookie")
	if err != nil {
		id, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "User-cookie",
			Value: id.String(),
		}
		http.SetCookie(w, c)
	}
	var u user
	if un, ok := dbSession[c.Value]; ok {
		u = dbUser[un]
	}
	return u
}
func alredylog(r *http.Request) bool {
	c, err := r.Cookie("User-cookie")
	if err != nil {
		return false
	}
	un := dbSession[c.Value]
	_, ok := dbUser[un]
	if !ok {
		return false
	}
	return true
}
