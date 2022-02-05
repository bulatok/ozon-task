package in_memory

import (
	"fmt"
	ss "github.com/patrickmn/go-cache"
	"time"
)

type InMemory struct{
	Timeout time.Duration
	Cache   *ss.Cache
}

func (im *InMemory) Open() error{
	return nil
}

func (im *InMemory) Close() error{
	return nil
}

func (im *InMemory) FindByURL(newLink string) (string, error){
	val, found := im.Cache.Get(newLink)
	if found == false{
		return "-1", fmt.Errorf("no such URL")
	}
	return val.(string), nil
}

func (im *InMemory) AddLink(newLink string, oldLink string) error{
	im.Cache.Set(newLink, oldLink, im.Timeout)
	return nil
}
