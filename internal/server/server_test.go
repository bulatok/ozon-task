package server

import (
	"bytes"
	"github.com/bulatok/ozon-task/internal/store"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestServer_hanldeMainGET(t *testing.T) {
	tests := []struct {
		name 		string
		urlRequest  string
		expectCode  int
		expectJSON  string
	}{
		{
			"bad url",
			"http://localhost:8080/?uuuuu=1",
			http.StatusBadRequest,
			`{"result":"no such URL"}`,
		},
		{
			"bad url",
			"http://localhost:8080/YH1foQ4FFKHepep",
			http.StatusBadRequest,
			`{"result":"no such URL"}`,
		},
		{
			"correct url",
			"http://localhost:8080/RFCzyRcRNU",
			http.StatusOK,
			`{"result":"https://ya.ru/"}`,
		},
	}

	serverTest := NewServer(store.CreateTEST())
	if err := serverTest.Store.Open(); err != nil{
		log.Fatal(err)
	}
	defer serverTest.Store.Close()

	store.CleanUp(serverTest.Store)
	store.AddUrl("https://ya.ru/", "RFCzyRcRNU", serverTest.Store) // for third test

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.urlRequest, b)


			serverTest.ServeHTTP(w, req)
			serverTest.hanldeMain()


			assert.Equal(t, tt.expectCode, w.Code)
			assert.Equal(t, tt.expectJSON, w.Body.String())
		})
	}
}


func TestServer_hanldeMainPOST(t *testing.T) {
	tests := []struct {
		name 		string
		vls 		url.Values
		expectCode  int
		expectJSON  string
	}{
		{
			"incorrect form value",
			url.Values{
				"blabla" : {"https://www.google.com/"},
			},
			http.StatusBadRequest,
			`{"result":"'' is incorrect URL"}`,

		},
		{
			"incorrect url",
			url.Values{
				"url" : {"google.com"},
			},
			http.StatusBadRequest,
			`{"result":"'google.com' is incorrect URL"}`,
		},
		{
			"correct url",
			url.Values{
				"url" : {"https://ya.ru/"},
			},
			http.StatusOK,
			`{"result":"http://localhost:8080/RFCzyRcRNU"}`,
		},
	}

	serverTest := NewServer(store.CreateTEST())
	if err := serverTest.Store.Open(); err != nil{
		log.Fatal(err)
	}
	defer serverTest.Store.Close()

	store.CleanUp(serverTest.Store)
	store.AddUrl("https://ya.ru/", "RFCzyRcRNU", serverTest.Store) // for third test

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(tt.vls.Encode()))
			if err != nil{
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


			serverTest.ServeHTTP(w, req)
			serverTest.hanldeMain()

			assert.Equal(t, tt.expectCode, w.Code)
			assert.Equal(t, tt.expectJSON, w.Body.String())
		})
	}
}