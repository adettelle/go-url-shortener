package storage

import (
	"time"

	"golang.org/x/exp/rand"
)

type PathStorage struct {
	Paths map[string]string
}

func New() *PathStorage {
	paths := make(map[string]string)
	return &PathStorage{Paths: paths}
}

func (p *PathStorage) GetPath(name string) string {
	return p.Paths[name]
}

func (p *PathStorage) AddPath(fullPath string) string {
	rangeStart := 5
	rangeEnd := 10
	offset := rangeEnd - rangeStart
	randLength := seededRand.Intn(offset) + rangeStart

	charSet := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	randString := stringWithCharset(randLength, charSet)

	p.Paths[randString] = fullPath

	return randString
}

var seededRand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)] // letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
