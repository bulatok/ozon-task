package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)


// Resp writes to header the status and optionally
// returns an error as JSON
func Resp(w *http.ResponseWriter, in string, status int, typeRequset string){
	(*w).WriteHeader(status)
	res := GiveJSON(in)
	log.Printf("%s with response %s", typeRequset, string(res))
	(*w).Write(res)
}

// IsValidUrl is URL validator
func IsValidUrl(in string) error {
	_, err := url.ParseRequestURI(in) // ???
	if err != nil {
		return fmt.Errorf("'%v' is incorrect URL", in)
	}
	return nil
}

// GiveJSON returns JSON as []byte
func GiveJSON(in string) []byte {
	type U struct {
		Result string `json:"result"`
	}
	tt := U{Result: in}
	res, _ := json.Marshal(tt)
	return res
}

func HandleReq(r io.Reader) (string, error){
	type coming struct{
		Url string `json:"url"`
	}
	d := json.NewDecoder(r)
	income := coming{Url: "-1"} // by default it will be "-1" to check if json request has {"url":"someURL"}
	if err := d.Decode(&income); err != nil{
		return "", err
	}
	if income.Url == "-1"{
		return "", fmt.Errorf("json request is incorrect")
	}
	return income.Url, nil
}
