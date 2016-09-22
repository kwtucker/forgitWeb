package lib

import (
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/system"
	"net/http"
)

// Logout ...
type Logout struct {
	Env  system.Application
	Sess *sessions.CookieStore
}

// Logout ...
func (l *Logout) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Grab the Session
	session, err := l.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.Values["authed"] = 0
	session.Save(r, w)
	http.Redirect(w, r, "http://"+l.Env.Config.HostString()+"/", http.StatusFound)
}
