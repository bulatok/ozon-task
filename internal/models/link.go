package models

import (
	"fmt"
	"hash/fnv"
	"math"
)

type Link struct{
	OriginalURL string `json:"origin_url"`
	ParsedURL   string `json:"new_url"`
}

// GetHash calculate the hash of string
func GetHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (l *Link) SetOriginalURL(originalURL string){
	l.OriginalURL = originalURL
}

// CreateNewURL:
// 1) Calculate a hash of the original string
// 2) Takes the remainder of the division by 63^10
// 3) So we got hash(originURL)%(63^10)_10 -> hash_63 number system
func (l *Link) SetNewURL() error{
	if l.OriginalURL == ""{
		return fmt.Errorf("need to set original url")
	}

	allowed := make(map[int]int32)
	idx := 0
	for i := 'a'; i <= 'z'; i += 1 {
		allowed[idx] = i
		idx += 1
	}
	for i := 'A'; i <= 'Z'; i += 1 {
		allowed[idx] = i
		idx += 1
	}
	for i := '0'; i <= '9'; i += 1 {
		allowed[idx] = i
		idx += 1
	}
	allowed[idx] = '_'

	mx := GetHash(l.OriginalURL)
	mx %= uint64(math.Pow(float64(len(allowed)), 10))

	newURL := ""
	nwBase := uint64(len(allowed))

	// mx(base=10) -> mxNew(base=len(allowed))
	for {
		if mx < nwBase {
			newURL += string(allowed[int(mx)])
			break
		} else if mx >= nwBase { // for illustration
			newURL += string(allowed[int(mx%nwBase)])
			mx /= uint64(len(allowed))
		}
	}
	l.ParsedURL = newURL
	return nil
}

