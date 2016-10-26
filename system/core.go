package system

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
)

/*
Application ...
Template is a pointer to the parsed templates
Config is a pointer to
*/
type Application struct {
	Template *template.Template
	Config   *Configuration
	AuthConf *oauth2.Config
}

// Init grabs all the config information
func (application *Application) Init(filename *string) {
	// Configuration available because it is a uppercase struct
	// and doesn't need to be imported as long as it is in the same package
	var config = &Configuration{}
	config.Init(filename)
	// so config can be accessed through application
	application.Config = config
	application.LoadTemplates()
	application.AuthConfig()
}

// LoadTemplates will grab all the templates
func (application *Application) LoadTemplates() error {
	t := template.New("base")
	pattern := filepath.Join(application.Config.TemplateDir, "*.html")
	application.Template = template.Must(t.ParseGlob(pattern))
	return nil
}

//AuthConfig ...
func (application *Application) AuthConfig() error {
	application.AuthConf = &oauth2.Config{
		ClientID:     application.Config.GithubClientID,
		ClientSecret: application.Config.GithubClientSecret,
		Scopes:       []string{"user", "repo"},
		Endpoint:     github.Endpoint,
	}
	return nil
}

// Route will take in the controller and the method within the controller.
func (application *Application) Route(controller interface{}, controllerMethodFunc string) httprouter.Handle {
	// handling the route
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// Since you can't see the values of the incoming interface, Reflect will get the value,
		// then you MethodByName to check if that interface has a method of the route string.
		// Then turn the returned value for the MethodByName to a Interface. Use type assertion to
		// convert the interface to the required type of method.
		methodValue := reflect.ValueOf(controller).MethodByName(controllerMethodFunc).Interface()
		method := methodValue.(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (map[string]interface{}, int))

		// Method returns two values. interface(data)= data about that page
		// and int(code)= status code
		data, code := method(w, r, ps)
		if data != nil {
			// Parse template and pass in data values
			contenttpl := Parse(application.Template, data["ContentTemplate"].(string), data)
			data["Content"] = template.HTML(contenttpl)
			tpl := Parse(application.Template, "base", data)

			// if the coded is okay render,
			// if code in not okay redirect to 404
			switch code {
			case http.StatusOK:
				w.Header().Set("Content-Type", "text/html")
				io.WriteString(w, tpl)
				return
			case http.StatusSeeOther, http.StatusFound:
				http.Redirect(w, r, tpl, code)
			}
		}
		http.Redirect(w, r, "/", code)
	}
	return fn
}

// NoViewRoute will take in the controller and the method within the controller.
func (application *Application) NoViewRoute(controller interface{}, controllerMethodFunc string) httprouter.Handle {
	// handling the route
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		methodValue := reflect.ValueOf(controller).MethodByName(controllerMethodFunc).Interface()
		method := methodValue.(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params))
		method(w, r, ps)
	}
	return fn
}
