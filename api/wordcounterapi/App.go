package wordcounterapi

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/interfaces"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/wordcounterapi/structs"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
	"github.com/gorilla/mux"
)

type App struct {
	Router      *mux.Router
	wordCounter interfaces.WordCounterInterface
}

var handlers []handler = []handler{
	{
		Path:    "/count/{text}",
		Handler: createHandler,
		Method:  "GET",
	},
}

func (app *App) Initialize(wordCounter interfaces.WordCounterInterface) {
	app.wordCounter = wordCounter
	app.Router = mux.NewRouter()

	for _, handler := range handlers {
		handerFunc := handler.Handler(app)
		app.Router.HandleFunc(handler.Path, handerFunc).Methods(handler.Method)
	}
}

func (app *App) Kill(status int) {
	os.Exit(status)
}

func createHandler(app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var resp structs.WordCounterResponse

		if text := mux.Vars(r)["text"]; util.IsEmptyString(text) {
			resp = structs.WordCounterResponse{
				Status: http.StatusBadRequest,
			}
		} else {
			resp = structs.WordCounterResponse{
				Status:         http.StatusOK,
				WordCollection: app.wordCounter.CountWords(text),
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
