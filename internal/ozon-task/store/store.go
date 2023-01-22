package store

import (
	"github.com/bulatok/ozon-task/internal/ozon-task/models"
	_ "github.com/lib/pq"
)

//go:generate mockgen
type LinksRepo interface {
	Name() string
	Close() error
	Save(link *models.Link) error
	Get(shortLink string) (*models.Link, error)
	Delete(originalLink string) error
}
