package server

import (
	"encoding/json"
	"fmt"
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
			type coming struct{
				Url string `json:"url"`
			}
			d := json.NewDecoder(r.Body)
			income := &coming{Url: "-1"} // by default it will be "-1" to check if json request has {"url":"someURL"}
			if err := d.Decode(income); err != nil{
				Resp(&w, fmt.Errorf("json query is incorrect").Error(), http.StatusBadRequest, "POST")
				return
			}
			if income.Url == "-1"{
				Resp(&w, fmt.Errorf("json query is incorrect").Error(), http.StatusBadRequest, "POST")
				return
			}


			urlPOST := income.Url
			if err := IsValidUrl(&urlPOST); err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest, "POST")
				return
			}


			urlResponse, err := CreateNewURL(urlPOST, s.Store) // here
			if err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest, "POST")
				return
			}
			Resp(&w, "http://localhost:8080/" + urlResponse, http.StatusOK, "POST")
		case http.MethodGet:
			urlGET := r.URL.String()[1:]
			if strings.Contains(r.URL.String(), "http://localhost:8080/") {
				urlGET = strings.Split(r.URL.String(), "http://localhost:8080/")[1]
			}


			urlResponse, err := FindURL(urlGET, s.Store) // here
			if err != nil {
				Resp(&w, err.Error(), http.StatusBadRequest, "GET")
				return
			}
			Resp(&w, urlResponse, http.StatusOK, "GET")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}