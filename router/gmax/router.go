package gmax

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	S APIServer
)

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

type APIServer struct {
	r    *mux.Router
	srv  *http.Server
	tmpl *template.Template
}

func (s *APIServer) StartServer() {
	var err error

	// serve static files
	// This will serve files under http://localhost:8000/static/<filename>
	// hardcoded for security reasons
	s.r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/assets"))))

	// templates
	// template.Must throw an error and exit
	// tpl = template.Must(template.ParseGlob("templates/*.html"))
	// use template tmpl.Execute(w, data) and useage in template: {{.PageTitle}} would get value from data.PageTitle
	s.tmpl, err = template.ParseGlob("web/templates/*")

	if err != nil {
		log.Panic("Error: ", err)
	}
	log.Println(s.tmpl.DefinedTemplates())

	// Handlers
	s.AddHandlers()

	// make connection string
	s_addr := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	log.Println("Listen on: %s", s_addr)

	s.srv = &http.Server{
		Handler: S.r,
		Addr:    s_addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(viper.GetInt("server.timeout.write")) * time.Second,
		ReadTimeout:  time.Duration(viper.GetInt("server.timeout.read")) * time.Second,
	}

	log.Fatal(s.srv.ListenAndServe())
}

// Add handlers
func (s *APIServer) AddHandlers() {
	// Handlers
	s.r.HandleFunc("/api/health", s.apiHealth)
	s.r.HandleFunc("/", s.indexHandler)

	// Custom. TODO: add here something from parameters
	s.r.HandleFunc("/message/", s.messageHandler)
}

func (s *APIServer) apiHealth(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (s *APIServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Index handler")

	if err := s.tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *APIServer) messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("message handler")

	if err := s.tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
