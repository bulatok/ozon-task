package models

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
	"time"
)

var (
	ErrEmptyOrigin = errors.New("origin must not be empty")
)

const (
	allowed = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_`
)

type Link struct {
	Original  string    `json:"origin_url,omitempty"`
	Short     string    `json:"short_url,omitempty"`
	CreatedAt time.Time `json:"-"`
	hash      string
}

func (l *Link) SetShortLink(servicePublicUrl string) error {
	if servicePublicUrl == "" {
		return ErrEmptyOrigin
	}

	linkHash, err := getLinkHash(l.Original)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(servicePublicUrl, "/") {
		servicePublicUrl += "/"
	}

	(*l).Short = fmt.Sprintf("%s%s", servicePublicUrl, linkHash)
	(*l).hash = linkHash
	return nil
}

// getLinkHash
//
// 1) Calculate a hash of the original string
//
// 2) Takes the remainder of the division by 63^10
//
// 3) So we got hash(originURL)%(63^10)_10 -> hash_63 number system
func getLinkHash(original string) (string, error) {
	if original == "" {
		return "", ErrEmptyOrigin
	}

	reminderByte := "0"

	getHash := func(s string) (uint64, error) {
		s += reminderByte
		h := fnv.New64a()
		_, err := h.Write([]byte(s))
		return h.Sum64(), err
	}

	mx, err := getHash(original)
	if err != nil {
		return "", err
	}
	mx %= uint64(math.Pow(float64(len(allowed)), 10))

	hashResult := ""
	nwBase := uint64(len(allowed))

	// mx(base=10) -> mxNew(base=len(allowed))
	for {
		if mx < nwBase {
			hashResult += string(allowed[int(mx)])
			break
		} else if mx >= nwBase { // for illustration
			hashResult += string(allowed[int(mx%nwBase)])
			mx /= uint64(len(allowed))
		}
	}
	return hashResult, nil
}

// GetUnderlineHash returns the original link hash
func (l *Link) GetUnderlineHash() string {
	return l.hash
}
