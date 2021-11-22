package store

import (
	"database/sql"
	"fmt"
	"time"
)

func AddUrl(origin_url, parsed_url string, s *Store) error{
	if _, found := FindByParsedURL(parsed_url, s); found == nil{
		return nil
	}
	sqlSteatment := `INSERT INTO urls (origin_url, parsed_url) VALUES($1, $2)` // urls is a table
	if _, err := s.db.Exec(sqlSteatment, origin_url, parsed_url); err != nil{
		return fmt.Errorf("Add url : %v", err)
	}
	return nil
}

func FindByParsedURL(parsed_url string, s *Store) (string, error) {
	var origin_url string
	if err := s.db.QueryRow("SELECT origin_url FROM urls WHERE parsed_url = $1",
		parsed_url).Scan(&origin_url); err != nil {
		if err == sql.ErrNoRows {
			return "-1", fmt.Errorf(`cannot find '%s': not found`, parsed_url)
		}
		return "-1", fmt.Errorf(`cannot find '%s': %w`, parsed_url, err)
	}
	return origin_url, nil
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