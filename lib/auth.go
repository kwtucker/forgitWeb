package lib

import (
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
		session.Values["authed"] = 1
		session.Values["code"] = code
		session.Save(r, w)
		http.Redirect(w, r, "/terminal/", http.StatusFound)
	}
}

// SessionCheck ...
func (a *Auth) SessionCheck(w http.ResponseWriter, r *http.Request) *sessions.Session {
	// Grab the Session
	session, err := a.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// If the person is not authed send them back to home
	if session.Values["authed"] != 1 {
		session.Values["authed"] = 0
		session.Save(r, w)
		http.Redirect(w, r, "http://"+a.Env.Config.HostString()+"/", http.StatusFound)
	}
	return session
}

// Logout will handle the event of clicking the logout button.
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Grab the Session
	session, err := a.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// set auth to 0 for not authed and save it. Then send user home.
	session.Values["authed"] = 0
	session.Save(r, w)
	http.Redirect(w, r, "http://"+a.Env.Config.HostString()+"/", http.StatusFound)
}
