package cache

import (
	"fmt"
	"sync"

	"github.com/bulatok/ozon-task/internal/ozon-task/models"
)

// formatKey is just for example, in fact the key could just be
// the originalLink
func formatKeyLink(short string) string {
	return fmt.Sprintf("link:%s", short)
}

type links struct {
	mx   sync.Mutex
	data map[string]interface{}
}

func (l *links) Save(link *models.Link) error {
	l.mx.Lock()
	l.data[formatKeyLink(link.Short)] = link
	l.mx.Unlock()
	return nil
}
func (l *links) Get(shortLink string) (*models.Link, error) {
	l.mx.Lock()
	defer l.mx.Unlock()
	d, exist := l.data[formatKeyLink(shortLink)]
	if !exist {
		return nil, models.ErrNotFound
	}

	dst, ok := d.(*models.Link)
	if !ok {
		return nil, models.ErrNotFound
	}
	return dst, nil
}

func (l *links) Delete(originalLink string) error {
	l.mx.Lock()
	delete(l.data, formatKeyLink(originalLink))
	l.mx.Unlock()
	return nil
}
