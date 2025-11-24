package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *APIConfig) newTemplateData() templateData {
	return templateData{}
}

func (app *APIConfig) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *APIConfig) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.TemplateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
