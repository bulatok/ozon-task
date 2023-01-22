package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bulatok/ozon-task/internal/ozon-task/models"
	"github.com/go-redis/redis/v8"
)

const (
	DefaultExpDur = time.Hour * 24 * 365
)

type links struct {
	rdb *redis.Client
}

// formatKey is just for example, in fact the key could just be
// the originalLink
func formatKeyLink(short string) string {
	return fmt.Sprintf("link:%s", short)
}

func decodeValueLink(link *models.Link) ([]byte, error) {
	link.CreatedAt = time.Now()
	return json.Marshal(link)
}

func encodeValueLink(data []byte) (*models.Link, error) {
	dst := &models.Link{}
	if err := json.Unmarshal(data, dst); err != nil {
		return nil, err
	}
	return dst, nil
}

func (l *links) Save(link *models.Link) error {
	key := formatKeyLink(link.Short)
	value, err := decodeValueLink(link)
	if err != nil {
		return err
	}

	return l.rdb.Set(l.rdb.Context(), key, value, DefaultExpDur).Err()
}

func (l *links) Get(shortLink string) (*models.Link, error) {
	key := formatKeyLink(shortLink)

	res, err := l.rdb.Get(l.rdb.Context(), key).Result()
	if err == redis.Nil {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return encodeValueLink([]byte(res))
}

func (l *links) Delete(shortLink string) error {
	key := formatKeyLink(shortLink)
	return l.rdb.Del(l.rdb.Context(), key).Err()
}
