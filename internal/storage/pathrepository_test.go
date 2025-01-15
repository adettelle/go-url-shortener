package storage

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddPath(t *testing.T) {
	pathStorage := New()

	short1, err := pathStorage.AddPath("http://localhost:8080/")
	require.NoError(t, err)
	for elem := range short1 {
		if elem >= 'a' && elem <= 'z' || elem >= 'A' && elem <= 'Z' {
			fullPath1, err := pathStorage.GetPath(short1)
			require.NoError(t, err)
			require.Equal(t, "http://localhost:8080/", fullPath1)
		} else {
			// require.Equal()
			log.Println("Error in char set") // TODO!!!!!!!!!!!
		}
	}

	short2, err := pathStorage.AddPath("http://localhost:8080/")
	require.NoError(t, err)
	for elem := range short2 {
		if elem >= 'a' && elem <= 'z' || elem >= 'A' && elem <= 'Z' {
			fullPath2, err := pathStorage.GetPath(short2)
			require.NoError(t, err)
			require.Equal(t, "http://localhost:8080/", fullPath2)
		} else {
			log.Println("Error in char set") // TODO!!!!!!!!!!!
		}
	}

}

func TestAddPathEmptyString(t *testing.T) {
	pathStorage := New()
	myErr := &EmptyPathError{}
	_, err := pathStorage.AddPath("")
	require.Equal(t, err, myErr)
}

func TestGetPath(t *testing.T) {
	pathStorage := New()

	fullPath := "http://localhost:8080/"
	name, err := pathStorage.AddPath(fullPath)
	require.NoError(t, err)

	short, err := pathStorage.GetPath(name)
	require.NoError(t, err)
	require.Equal(t, fullPath, short)
}

func TestGetPathUnknownName(t *testing.T) {
	pathStorage := New()

	unknownName := "aaa"
	_, err := pathStorage.GetPath(unknownName)
	require.Equal(t, err, &NoEntryError{name: unknownName})
}

func TestStringWithCharset(t *testing.T) {
	charSet := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

	newStr, err := stringWithCharset(10, charSet)
	require.NoError(t, err)
	require.Len(t, newStr, 10)
}

func TestStringWithCharsetInvalidLength(t *testing.T) {
	charSet := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

	newStr, err := stringWithCharset(-10, charSet)
	require.Equal(t, err, &InvalidLengthError{})
	require.Empty(t, newStr)
}

func TestStringWithCharsetInvalidCharSet(t *testing.T) {
	charSet := ""

	newStr, err := stringWithCharset(10, charSet)
	require.Equal(t, err, &InvalidCharSetError{})
	require.Empty(t, newStr)
}
