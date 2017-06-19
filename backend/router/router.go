package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ferreus/jobs/backend/scrapper"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	/*for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Authorize(handler, route.Name, route.Role)
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}*/
	var handler http.HandlerFunc
	handler = Jobs
	router.Methods("GET").Path("/v1/jobs").Name("Jobs").Handler(handler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend/dist/")))

	return router
}

func Jobs(w http.ResponseWriter, r *http.Request) {
	url := "https://www.drushim.co.il/jobs/cat6/?area=3"
	var jobItems []scrapper.JobItem
	log.Print("Fetching page:0...\n")
	items, next, err := scrapper.Fetch(url)
	if err != nil {
		log.Panic(err)
	}
	jobItems = items
	for i := 0; i < 10; i++ {
		log.Printf("Fetching page:%d [%s]...\n", i+1, next)
		items, next, err = scrapper.Fetch(next)
		if err != nil {
			log.Panic(err)
		}
		jobItems = append(jobItems, items...)
	}
	jsonTxt, err := json.Marshal(jobItems)
	if err != nil {
		log.Panic(err)
	}
	log.Print(string(jsonTxt))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(jobItems); err != nil {
		log.Panic(err)
	}
}
