package server

import (
	"github.com/bulatok/ozon-task/configs"
	"github.com/bulatok/ozon-task/internal/models"
	"github.com/bulatok/ozon-task/internal/store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

// Start является точкой входа сервера
func Start(config *configs.Config, store store.Store) error {
	log.Println("Start listening on port :8080")

	srv := NewServer(store)
	return http.ListenAndServe(config.Port, srv)
}

type Server struct{
	Router *mux.Router
	Store  store.Store
}
func NewServer(str store.Store) *Server{
	s :=  &Server{
		Router: mux.NewRouter(),
		Store:  str,
	}

	s.Router.PathPrefix("/").Handler(s.hanldeMain())
	return s
}


func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request){
	s.Router.ServeHTTP(w, r)
}


func (s *Server) hanldeMain() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			oldURL, err := HandleReq(r.Body)
			if err != nil{
				Resp(&w, err.Error(), http.StatusBadRequest, "POST")
			}

			l := models.Link{}
			l.SetOriginalURL(oldURL)

			if err := IsValidUrl(l.OriginalURL); err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest, "POST")
				return
			}

			if err := l.SetNewURL(); err != nil{
				Resp(&w, err.Error(), http.StatusBadRequest, "POST")
				return
			}

			if err := s.Store.AddLink(l.ParsedURL, l.OriginalURL); err != nil{
				Resp(&w, err.Error(), http.StatusBadRequest, "POST")
				return
			}

			Resp(&w, "http://localhost:8080/" + l.ParsedURL, http.StatusOK, "POST")
		case http.MethodGet:
			newURL := r.URL.String()[1:]
			if strings.Contains(r.URL.String(), "http://localhost:8080/") {
				newURL = strings.Split(r.URL.String(), "http://localhost:8080/")[1]
			}

			oldURL, err := s.Store.FindByURL(newURL)
			if err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest, "GET")
				return
			}
			Resp(&w, oldURL, http.StatusOK, "GET")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}