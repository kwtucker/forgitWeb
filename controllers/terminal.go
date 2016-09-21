package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit_noFrame/models"
	"github.com/kwtucker/forgit_noFrame/system"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

// TerminalController ...
type TerminalController struct {
	Env  system.Application
	sess *sessions.CookieStore
}

// Terminal ...
func (c *TerminalController) Terminal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	// Grab the Session
	session, err := c.sess.Get(r, "ForgitSession")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Get values from URL
	code := r.FormValue("code")
	state := r.FormValue("state")

	// Check if github returned the same state that I made the request with.
	if state != c.Env.Config.GithubState {
		session.Values["authed"] = 0
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	} else {
		// logged in
		session.Values["authed"] = 1
	}
}
