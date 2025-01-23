package storage

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

const charSet = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

type AddressStorage struct {
	Addresses map[string]string
}

func New() *AddressStorage {
	addresses := make(map[string]string)
	return &AddressStorage{Addresses: addresses}
}

type NoEntryError struct {
	name string
}

func (e *NoEntryError) Error() string {
	return fmt.Sprintf("No Entry for name %s", e.name)
}

// возращает полный url по ключу (короткому url)
func (a *AddressStorage) GetAddress(name string) (string, error) {
	if addr, ok := a.Addresses[name]; ok {
		return addr, nil
	}

	return "", &NoEntryError{
		name: name,
	}
}

type EmptyAddressError struct{}

func (e *EmptyAddressError) Error() string {
	return "Empty full address"
}

func (a *AddressStorage) AddAddress(fullAddress string) (string, error) {
	if fullAddress == "" {
		return "", &EmptyAddressError{}
	}
	rangeStart := 2
	rangeEnd := 10
	offset := rangeEnd - rangeStart
	randLength := seededRand.Intn(offset) + rangeStart

	randString, err := stringWithCharset(randLength, charSet)
	if err != nil {
		return "", err
	}

	a.Addresses[randString] = fullAddress

	return randString, nil
}

type InvalidLengthError struct{}

func (e *InvalidLengthError) Error() string {
	return "Invalid length!"
}

type InvalidCharSetError struct{}

func (e *InvalidCharSetError) Error() string {
	return "Invalid charset!"
}

var seededRand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

// stringWithCharset generates a random string of a specified length using the provided character set.
// Parameters:
//   - length: The desired length of the generated string.
//   - charset: A string containing the characters to use for generating the random string.
//
// Returns:
//   - A random string composed of characters from the given charset.
func stringWithCharset(length int, charset string) (string, error) {
	if length <= 0 {
		return "", &InvalidLengthError{}
	}

	if charset == "" {
		return "", &InvalidCharSetError{}
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}
	return string(b), nil
}
