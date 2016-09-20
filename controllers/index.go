package controllers

import (
	// "fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// IndexController ...
type IndexController struct {
}

// Index ...
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int) {

	// Nav for this view.
	navLinks := map[string]string{
		"#":          "Features",
		"#pricing":   "Pricing",
		"#createdby": "Created By",
	}
	// values for the view.
	data := map[string]interface{}{
		"PageName":        "Forgit",
		"ContentTemplate": "index",
		"NavLinks":        navLinks,
	}
	return data, http.StatusOK
}
