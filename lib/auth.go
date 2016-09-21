package lib

import (
	// "fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit_noFrame/system"
	"net/http"
)

// Auth is a struct that gives access to the application
type Auth struct {
	Env system.Application
}

// AuthFunc will redirect to github
func (a *Auth) AuthFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginString := a.Env.AuthConf.AuthCodeURL(a.Env.Config.GithubState)
	http.Redirect(w, r, loginString, http.StatusTemporaryRedirect)
}
