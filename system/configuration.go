package system

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Configuration ...
type Configuration struct {
	SessionSecret             string
	GithubClientID            string
	GithubClientSecret        string
	GithubState               string
	WebPort                   int
	DbPort                    int
	DbName                    string
	TemplateDir               string
	TemplatePreCompile        bool
	Debug                     bool
	GoogleAnalyticsTrackingID string
	WebHost                   string
	DbHost                    string
	Streams                   []map[string]string
	StaticPath                string
}

// Init reads the config file and sets all values to config struct
func (c *Configuration) Init(filename *string) {

	if ecp := os.Getenv("FORGIT_CONFIG_PATH"); ecp != "" {
		filename = &ecp
	}

	file, err := os.Open(*filename)
	if err != nil {
		if len(*filename) > 1 {
			fmt.Printf("Error: could not read config file %s.\n", *filename)
		}
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	// Overwrite the defaults
	if err := decoder.Decode(&c); err == io.EOF {
		fmt.Println(err)
	} else if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.String())

}

// HostString ...
func (c *Configuration) HostString() string {
	return fmt.Sprintf("%s:%d", c.WebHost, c.WebPort)
}

// WebPortString returns a port string
func (c *Configuration) WebPortString() string {
	return fmt.Sprintf(":%d", c.WebPort)
}

// DbHostString returns Mongo Connection string
func (c *Configuration) DbHostString() string {
	if c.DbPort > 0 {
		return fmt.Sprintf("mongodb://%s:%d", c.DbHost, c.DbPort)
	}
	return fmt.Sprintf("mongodb://%s", c.DbHost)
}

func (c *Configuration) String() string {
	s := "Config:\n"
	s += fmt.Sprintf("   Host: %s,\n", c.HostString())
	s += fmt.Sprintf("   DB: %s,\n", c.DbHostString())
	s += fmt.Sprintf("   TemplatePath: %s,\n", c.TemplateDir)
	s += fmt.Sprintf("   StaticPath: %s,\n", c.StaticPath)
	s += fmt.Sprintf("   TemplatePreCompile: %v,\n", c.TemplatePreCompile)
	s += fmt.Sprintf("   Debug: %v\n", c.Debug)
	s += fmt.Sprintf("   GoogleAnalyticsTrackingID: %v\n", c.GoogleAnalyticsTrackingID)
	return s
}
