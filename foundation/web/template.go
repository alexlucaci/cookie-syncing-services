package web

import (
	"github.com/pkg/errors"
	"html/template"
	"net/http"
)

// RenderTemplate renders a given template along with the data passed
// It will error out if parsing the template or executing it fails for some reason.
func RenderTemplate(w http.ResponseWriter, filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return errors.Wrapf(err, "parsing template file")
	}

	err = t.Execute(w, data)
	if err != nil {
		return errors.Wrapf(err, "executing template")
	}

	return nil
}
