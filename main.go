package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserName string
	Password []byte
	FistName string
	LastName string
	Role     string
}

var tpl *template.Template
var dbSession = map[string]string{}
var dbUser = map[string]user{}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	fmt.Println("Starting...")
	http.HandleFunc("/", ind)
	http.HandleFunc("/sing", singup)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/log", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/Admin", adm)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
func ind(w http.ResponseWriter, r *http.Request) {
	u := getus(w, r)
	err := tpl.ExecuteTemplate(w, "index.gohtml", u)
	if err != nil {
		log.Fatal(err)
	}
}
func singup(w http.ResponseWriter, r *http.Request) {
	if alredylog(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	var u user
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")
		rol := r.FormValue("role")
		if _, ok := dbUser[un]; ok {
			http.Error(w, "This un is taken!", http.StatusForbidden)
			return
		}
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "User-cookie",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSession[c.Value] = un
		bc, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "intelnal serv err", http.StatusForbidden)
			return
		}
		u = user{un, bc, f, l, rol}
		dbUser[un] = u
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "singup.gohtml", u)
}
func bar(w http.ResponseWriter, r *http.Request) {
	u := getus(w, r)
	if !alredylog(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if u.Role == "Admin" {
		http.Error(w, "You must go to dmin panel", http.StatusForbidden)
	}
	err := tpl.ExecuteTemplate(w, "bar.gohtml", u)
	if err != nil {
		log.Fatal(err)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	if alredylog(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		u, ok := dbUser[un]
		if !ok {
			http.Error(w, "Username do not macth", http.StatusForbidden)
			return
		}
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Password do not macth", http.StatusForbidden)
			return
		}
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "User-cookie",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSession[c.Value] = un
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
func logout(w http.ResponseWriter, r *http.Request) {
	if !alredylog(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, _ := r.Cookie("User-cookie")
	delete(dbSession, c.Value)
	c = &http.Cookie{
		Name:   "User-cookie",
		Value:  " ",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/log", http.StatusSeeOther)
}
func adm(w http.ResponseWriter, r *http.Request) {
	if !alredylog(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var u user
	if u.Role != "Admin" {
		http.Error(w, "U are not admin!!!", http.StatusForbidden)
		return
	}
	us := getus(w, r)

	tpl.ExecuteTemplate(w, "admin.gohtml", us)
}
