package store

import (
	"fmt"
	"time"
)

func AddUrl(oldURL, newURL string, s *Store) {
	sqlSteatment := `INSERT INTO urls (origin_url, parsed_url) VALUES($1, $2)` // urls is a table
	s.db.QueryRow(sqlSteatment, oldURL, newURL)
}

func FindByParsedURL(parsed_url string, s *Store) (string, error) {
	sqlSteatment := fmt.Sprintf(`SELECT origin_url FROM urls WHERE parsed_url = '%s'`, parsed_url)
	rows, err := s.db.Query(sqlSteatment)
	if err != nil {
		return "-1", err
	}
	err = fmt.Errorf("no such URL") // default answer will be such like
	res :=  ""
	for rows.Next() {
		rows.Scan(&res)
		err = nil
	}
	return res, err
}
func CleanUp(s *Store){
	s.db.Exec(`DELETE FROM urls`)
	s.db.Exec("ALTER SEQUENCE urls_id_seq RESTART WITH 1")
}

// In memory
func FindByParsedURLInMemory(parsed_url string, s *Store) (string, error){
	val, found := s.cch.Get(parsed_url)
	if found == false{
		return "-1", fmt.Errorf("no such URL")
	}
	return val.(string), nil
}

func AddUrlInMemory(oldURL, newURL string, s *Store) {
	s.cch.Set(newURL, oldURL, time.Minute * 5) // for example 5 minute??
}