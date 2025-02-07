package storage

import "fmt"

type EmptyAddressError struct{}

func (e *EmptyAddressError) Error() string {
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
