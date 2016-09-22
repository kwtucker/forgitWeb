package controllers

import (
	// "encoding/json"
	"fmt"
	// "github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	// "github.com/kwtucker/forgit_noFrame/models"
	"github.com/kwtucker/forgit/system"
	// "golang.org/x/oauth2"
	// "io/ioutil"
	// "log"
	"net/http"
)

// TerminalController ...
type TerminalController struct {
	Env  system.Application
	Sess *sessions.CookieStore
}

// Terminal ...
func (c *TerminalController) Terminal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	// Grab the Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// If the person is not authed send them back to home
	if session.Values["authed"] != 1 {
		session.Values["authed"] = 0
		fmt.Println("authed not term", session.Values["authed"])
		session.Save(r, w)
		http.Redirect(w, r, "http://"+c.Env.Config.HostString()+"/", http.StatusFound)
	}
	// Nav for this view.
	navLinks := map[string]string{
		"/terminal/":      "Terminal",
		"#gettingstarted": "Getting Started",
		"/logout":         "Logout",
	}
	// data for the view
	data := map[string]interface{}{
		"PageName":        "Terminal",
		"ContentTemplate": "terminal",
		"NavLinks":        navLinks,
	}
	return data, http.StatusOK
}
