package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	sw "github.com/ferreus/jobs/backend/router"

	"github.com/gorilla/handlers"
)

func main() {
	port := 8080
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server started on " + addr)
	router := sw.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}
