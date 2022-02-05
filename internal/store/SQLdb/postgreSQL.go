package SQLdb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgreDB struct {
	DB    *sql.DB
	DbURL string
}

func (pdb *PostgreDB) Open() error{
	db, err := sql.Open("postgres", pdb.DbURL)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	pdb.DB = db
	return nil
}

func (pdb *PostgreDB) Close() error{
	if err := pdb.DB.Close(); err != nil{
		return err
	}
	return nil
}

func (pdb *PostgreDB) FindByURL(newLink string) (string, error){
	var oldLink string
	if err := pdb.DB.QueryRow("SELECT origin_url FROM urls WHERE parsed_url = $1",
		newLink).Scan(&oldLink); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf(`cannot find '%s': not found`, newLink)
		}
		return "", fmt.Errorf(`cannot find '%s': %w`, newLink, err)
	}
	return oldLink, nil
}

func (pdb *PostgreDB) AddLink(newLink string, oldLink string) error{

	// если ссылка уже лежит в бд
	if _, err := pdb.FindByURL(newLink); err == nil{
		return nil
	}

	if _, err := pdb.DB.Exec("INSERT INTO urls (origin_url, parsed_url) VALUES($1, $2)", oldLink, newLink); err != nil{
		return fmt.Errorf("Add url : %v", err)
	}
	return nil
}