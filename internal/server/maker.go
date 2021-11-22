package server

import (
	"encoding/json"
	"fmt"
	"github.com/bulatok/ozon-task/internal/store"
	"hash/fnv"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
)


type U struct {
	Result string `json:"result"`
}


func Resp(w *http.ResponseWriter, in string, status int, typeRequset string){
	(*w).WriteHeader(status)
	res := GiveJSON(in)
	log.Printf("%s with response %s", typeRequset, string(res))
	(*w).Write(res)
}


func IsValidUrl(in *string) error {
	_, err := url.ParseRequestURI(*in) // ???
	if err != nil {
		return fmt.Errorf("'%v' is incorrect URL", *in)
	}
	return nil
}


func GiveJSON(in string) []byte {
	tt := U{Result: in}
	res, _ := json.Marshal(tt)
	return res
}


func GetHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// GetCoded:
// 1) Calculate a hash of the original string
// 2) Takes the remainder of the division by 63^10
// 3) So we got hash(originURL)%(63^10)_10 -> hash_63 number system
func GetCoded(originURL string) string {
	allowed := make(map[int]int32)
	idx := 0
	for i := 'a'; i <= 'z'; i += 1 {
		allowed[idx] = i
		idx += 1
	}
	for i := 'A'; i <= 'Z'; i += 1 {
		allowed[idx] = i
		idx += 1
	}
	for i := '0'; i <= '9'; i += 1 {
		allowed[idx] = i
		idx += 1
	}
	allowed[idx] = '_'

	mx := GetHash(originURL)
	mx %= uint64(math.Pow(float64(len(allowed)), 10))

	newURL := ""
	nwBase := uint64(len(allowed))

	// mx(base=10) -> mxNew(base=len(allowed))
	for {
		if mx < nwBase {
			newURL += string(allowed[int(mx)])
			break
		} else if mx >= nwBase { // for illustration
			newURL += string(allowed[int(mx%nwBase)])
			mx /= uint64(len(allowed))
		}
	}
	return newURL
}


func CreateNewURL(originURL string, s *store.Store) (string, error) {
	// https://?url=0123456789 ???
	newURL := GetCoded(originURL)
	switch s.TypeDB{
	case "Postgres":
		if err := store.AddUrl(originURL, newURL, s); err != nil{
			return "-1", err
		}
	default:
		store.AddUrlInMemory(originURL, newURL, s)
	}

	return newURL, nil
}


func FindURL(parsedURL string, s *store.Store) (string, error) {
	var res string
	var err error
	switch s.TypeDB{
	case "Postgres":
		res, err = store.FindByParsedURL(parsedURL, s)
	default:
		res, err = store.FindByParsedURLInMemory(parsedURL, s)
	}
	if err != nil {
		return "-1", err
	}
	return res, nil
}


func IsRequestOK(in string) error{
	// /?url=.... || /cut?url=...
	if strings.Contains(in, "/?url=") || strings.Contains(in, "/cut?url="){
		return nil
	}
	return fmt.Errorf("invalid request")
}
