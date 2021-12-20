package wordcounterapi

import "net/http"

type handler struct {
	Path    string
	Handler func(app *App) func(w http.ResponseWriter, r *http.Request)
	Method  string
}
