package storage

import (
	"fmt"
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

type NoEntryError struct {
	name string
}

func (e *NoEntryError) Error() string {
	return fmt.Sprintf("No Entry for name %s", e.name)
}

func (p *PathStorage) GetPath(name string) (string, error) {
	if p.Paths[name] == "" {
		return "", &NoEntryError{
			name: name,
		}
	}
	return p.Paths[name], nil
}

type EmptyPathError struct{}

func (e *EmptyPathError) Error() string {
	return "Empty full path!"
}

func (p *PathStorage) AddPath(fullPath string) (string, error) {
	if fullPath == "" {
		return "", &EmptyPathError{}
	}
	rangeStart := 2
	rangeEnd := 10
	offset := rangeEnd - rangeStart
	randLength := seededRand.Intn(offset) + rangeStart

	charSet := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	randString, err := stringWithCharset(randLength, charSet)
	if err != nil {
		return "", err
	}

	p.Paths[randString] = fullPath

	return randString, nil
}

type InvalidLengthError struct{}

func (e *InvalidLengthError) Error() string {
	return "Invalid length!"
}

type InvalidCharSetError struct{}

func (e *InvalidCharSetError) Error() string {
	return "Invalid char set!"
}

var seededRand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

// stringWithCharset generates a random string of a specified length using the provided character set.
// Parameters:
//   - length: The desired length of the generated string.
//   - charset: A string containing the characters to use for generating the random string.
//
// Returns:
//   - A random string composed of characters from the given charset.
//
// Notes:
//   - The function uses a seeded random number generator to ensure randomness.
//   - If the charset is empty or the length is non-positive, the behavior is undefined. TODO!!!!!!!!!!!!!!!!!
func stringWithCharset(length int, charset string) (string, error) {
	if length <= 0 {
		return "", &InvalidLengthError{}
	}

	if charset == "" {
		return "", &InvalidCharSetError{}
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)] // letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b), nil
}
