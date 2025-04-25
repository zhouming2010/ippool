package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	gid := 3
	fmt.Fprintf(w, "Welcome to the Home Page! gid=%d\n", gid)
}

func Start() {
	r := mux.NewRouter()
	// var handler RouterHandler
	r.HandleFunc("/", mainHandler)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	log.Printf("Starting server on : %d...\n", GetAppConf().IPServer.Port)
	addr := fmt.Sprintf(":%d", GetAppConf().IPServer.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
