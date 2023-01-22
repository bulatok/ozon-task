package postgres

import (
	"database/sql"

	"github.com/bulatok/ozon-task/internal/ozon-task/models"

	_ "github.com/lib/pq"
)

type links struct {
	*sql.DB
}

func (l *links) exists(link *models.Link) (bool, error) {
	cnt := 0
	err := l.DB.
		QueryRow("SELECT count(*) FROM links WHERE original_url = $1 AND short_url = $2", link.Original, link.Short).
		Scan(&cnt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return cnt != 0, nil
}

func (l *links) Save(link *models.Link) error {
	exist, err := l.exists(link)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	_, err = l.DB.
		Exec("INSERT INTO links (original_url, short_url) VALUES ($1, $2)",
			link.Original, link.Short)
	return err
}

func (l *links) Get(shortLink string) (*models.Link, error) {
	originalLink := ""
	err := l.DB.
		QueryRow("SELECT original_url FROM links WHERE short_url = $1", shortLink).
		Scan(&originalLink)

	if err == sql.ErrNoRows {
		return nil, models.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &models.Link{
		Short:    shortLink,
		Original: originalLink,
	}, nil
}

func (l *links) Delete(originalLink string) error {
	_, err := l.DB.Exec("DELETE FROM links WHERE original_url = $1", originalLink)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
