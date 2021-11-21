package server

import (
	"github.com/bulatok/ozon-task/internal/store"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Server struct{
	Router *mux.Router
	Store *store.Store
}


func NewServer(str *store.Store) *Server{
	s :=  &Server{
		Router: mux.NewRouter(),
		Store: str,
	}

	s.Router.PathPrefix("/").Handler(s.hanldeMain())
	return s
}


func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request){
	s.Router.ServeHTTP(w, r)
}


func (s * Server) hanldeMain() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			r.ParseForm()
			urlPOST := r.FormValue("url")
			if err := IsValidUrl(&urlPOST); err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest)
				return
			}
			urlResponse, err := CreateNewURL(urlPOST, s.Store) // here
			if err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest)
				return
			}
			Resp(&w, "http://localhost:8080/" + urlResponse, http.StatusOK)
		case http.MethodGet:
			urlGET := r.URL.String()[1:]
			if strings.Contains(r.URL.String(), "http://localhost:8080/") {
				urlGET = strings.Split(r.URL.String(), "http://localhost:8080/")[1]
			}
			urlResponse, err := FindURL(urlGET, s.Store) // here
			if err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest)
				return
			}
			Resp(&w, urlResponse, http.StatusOK)


			// here I implemented redirecting to the page
			//urlGET := r.URL.String()[1:]
			//log.Println(urlGET)
			//urlResponse, err := FindURL(urlGET, s.Store)
			//if err != nil {
			//	Resp(&w, err.Error(), http.StatusBadRequest)
			//	return
			//}
			//http.Redirect(w, r, urlResponse, http.StatusSeeOther) // redirect to page
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}