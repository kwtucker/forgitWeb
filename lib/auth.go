package lib

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/system"
	"net/http"
)

// Auth is a struct that gives access to the application
type Auth struct {
	Env  system.Application
	Sess *sessions.CookieStore
}

// AuthFunc will redirect to github
func (a *Auth) AuthFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginString := a.Env.AuthConf.AuthCodeURL(a.Env.Config.GithubState)
	http.Redirect(w, r, loginString, http.StatusTemporaryRedirect)
}

// Callback ...
func (a *Auth) Callback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Grab the Session
	session, err := a.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Get values from URL
	code := r.FormValue("code")
	state := r.FormValue("state")
	if code == "" || state == "" || state != a.Env.Config.GithubState {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		// logged in
		fmt.Println("authed")
		session.Values["authed"] = 1
		session.Save(r, w)
		fmt.Println("authed", session.Values["authed"])
		fmt.Println("http://" + a.Env.Config.HostString() + "/terminal/")
		http.Redirect(w, r, "/terminal/", http.StatusFound)
	}
}
