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

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["sub_router"] // 获取路径参数

	fmt.Fprintf(w, "Request handled for ID: %s", id)
}

func handlSubRouter(w http.ResponseWriter, r *http.Request, handlers map[string]http.HandlerFunc) {
	vars := mux.Vars(r)
	id := vars["sub_router"]
	handler := handlers[id]
	if handler == nil {
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/user/{sub_router}", userHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.HandleFunc("/about", aboutHandler)

	log.Printf("Starting server on : %d...\n", GetAppConf().IPServer.Port)
	addr := fmt.Sprintf(":%d", GetAppConf().IPServer.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
