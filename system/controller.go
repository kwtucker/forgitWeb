package system

import (
	"bytes"
	"fmt"
	"html/template"
)

// Parse will build the html and pass all the data values
func Parse(t *template.Template, name string, data map[string]interface{}) string {
	// Create bytes buffer for io.writer param in the ExecuteTemplate.
	var doc bytes.Buffer
	err := t.ExecuteTemplate(&doc, name, data)
	if err != nil {
		fmt.Printf("template execution: %s", err)
	}
	return doc.String()
}
