package storage

import "fmt"

type EmptyOriginalURLError struct{}

func (e *EmptyOriginalURLError) Error() string {
	return "Empty full address"
}

type InvalidLengthError struct{}

func (e *InvalidLengthError) Error() string {
	return "Invalid length!"
}

type InvalidCharSetError struct{}

func (e *InvalidCharSetError) Error() string {
	return "Invalid charset!"
}

type NoEntryError struct {
	ShortURL string
}

func (e *NoEntryError) Error() string {
	return fmt.Sprintf("No Entry for shortURL %s", e.ShortURL)
}

type OriginalURLExistsErr struct {
	// ShortURL    string
	OriginalURL string
}

func (e *OriginalURLExistsErr) Error() string {
	return fmt.Sprintf("originalURL is already exists: %s", e.OriginalURL) // , e.ShortURL
}

func NewOriginalURLExistsErr(shortURL string, originalURL string) *OriginalURLExistsErr {
	return &OriginalURLExistsErr{
		// ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
