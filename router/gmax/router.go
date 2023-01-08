package gmax

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	S APIServer
)

type APIServer struct {
	r   *mux.Router
	srv *http.Server
}

func init() {
	// create router
	S.r = mux.NewRouter()

	// 	r.HandleFunc("/books/{title}", CreateBook).Methods("POST").Host("www.mybookstore.com")
	// r.HandleFunc("/books/{title}", ReadBook).Methods("GET").Schemes("https")
	// r.HandleFunc("/books/{title}", UpdateBook).Methods("PUT")
	// r.HandleFunc("/books/{title}", DeleteBook).Methods("DELETE")

	// 	bookrouter := r.PathPrefix("/books").Subrouter()
	// bookrouter.HandleFunc("/", AllBooks)
	// bookrouter.HandleFunc("/{title}", GetBook)

}

func (s *APIServer) StartServer() {
	s.r.HandleFunc("/api/health", apiHealth)

	// serve static files
	// This will serve files under http://localhost:8000/static/<filename>
	// hardcoded for security reasons
	s.r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/assets"))))

	// make connection string
	s_addr := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	log.Println("Listen on: %s", s_addr)

	s.srv = &http.Server{
		Handler: S.r,
		Addr:    s_addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(s.srv.ListenAndServe())
}

func apiHealth(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
