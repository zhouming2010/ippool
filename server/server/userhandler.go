package server

import (
	"fmt"
	"net/http"
)

var userHandlers = map[string]http.HandlerFunc{
	"login":    loginHandler,
	"register": registerHandler,
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	gid := 1
	fmt.Fprintf(w, "Welcome to the Home Page! gid=%d\n", gid)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	gid := 2
	fmt.Fprintf(w, "Welcome to the Home Page! gid=%d\n", gid)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	handlSubRouter(w, r, userHandlers)
}
