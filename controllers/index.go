package controllers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/kwtucker/forgit/system"
	"net/http"
)

// IndexController ...
type IndexController struct {
	Sess *sessions.CookieStore
	Env  system.Application
}

// Index ...
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {
	// Get Session
	session, err := c.Sess.Get(r, "ForgitSession")
	if err != nil {
		fmt.Println("damn index")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// if the session is new make authed to 0
	if session.IsNew {
		session.Values["authed"] = 0
		session.Save(r, w)
	}

	// values for the view.
	data := map[string]interface{}{
		"PageName":        "Forgit",
		"ContentTemplate": "index",
	}
	return data, http.StatusOK
}
