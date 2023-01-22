package cache

const (
	storeName = "cache"
)

type Cache struct {
	*links
}

func Provide() (*Cache, error) {
	return &Cache{
		&links{
			data: map[string]interface{}{},
		},
	}, nil
}

func (p *Cache) Close() error {
	return nil
}

func (p *Cache) Name() string {
	return storeName
}
