package storage

import (
	"testing"

	"github.com/adettelle/go-url-shortener/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var errlog *zap.Logger = logger.Logger

func TestAddAddress(t *testing.T) {
	addressStorage := New()

	short1, err := addressStorage.AddAddress("http://localhost:8080/")
	require.NoError(t, err)
	for elem := range short1 {
		if elem >= 'a' && elem <= 'z' || elem >= 'A' && elem <= 'Z' {
			fullAddress1, err := addressStorage.GetAddress(short1)
			require.NoError(t, err)
			require.Equal(t, "http://localhost:8080/", fullAddress1)
		} else {
			// require.Equal() ??? TODO
			errlog.Error("Error in charset", zap.Error(err))

		}
	}

	short2, err := addressStorage.AddAddress("http://localhost:8080/")
	require.NoError(t, err)
	for elem := range short2 {
		if elem >= 'a' && elem <= 'z' || elem >= 'A' && elem <= 'Z' {
			fullAddress2, err := addressStorage.GetAddress(short2)
			require.NoError(t, err)
			require.Equal(t, "http://localhost:8080/", fullAddress2)
		} else {
			// require.Equal() ??? TODO
			errlog.Error("Error in charset", zap.Error(err))
		}
	}

}

func TestAddAddressEmptyString(t *testing.T) {
	addressStorage := New()
	myErr := &EmptyAddressError{}
	_, err := addressStorage.AddAddress("")
	require.Equal(t, err, myErr)
}

func TestGetAddress(t *testing.T) {
	addressStorage := New()

	fullAddress := "http://localhost:8080/"
	name, err := addressStorage.AddAddress(fullAddress)
	require.NoError(t, err)

	short, err := addressStorage.GetAddress(name)
	require.NoError(t, err)
	require.Equal(t, fullAddress, short)
}

func TestGetAddressUnknownName(t *testing.T) {
	addressStorage := New()

	unknownName := "aaa"
	_, err := addressStorage.GetAddress(unknownName)
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
