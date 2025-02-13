package urlstorage

import (
	"log"
	"testing"

	"github.com/adettelle/go-url-shortener/internal/logger"
	"github.com/adettelle/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var errlog *zap.Logger = logger.Logger

func TestAddOriginalURL(t *testing.T) {
	addressStorage, err := New(false, "")
	if err != nil {
		log.Fatal(err)
	}

	short1, err := addressStorage.AddOriginalURL("http://localhost:8080/")
	require.NoError(t, err)
	for elem := range short1 {
		if elem >= 'a' && elem <= 'z' || elem >= 'A' && elem <= 'Z' {
			fullAddress1, err := addressStorage.GetOriginalURLByShortURL(short1)
			require.NoError(t, err)
			require.Equal(t, "http://localhost:8080/", fullAddress1)
		} else {
			// require.Equal() ??? TODO
			errlog.Error("Error in charset", zap.Error(err))

		}
	}

	short2, err := addressStorage.AddOriginalURL("http://localhost:8080/")
	require.NoError(t, err)
	for elem := range short2 {
		if elem >= 'a' && elem <= 'z' || elem >= 'A' && elem <= 'Z' {
			fullAddress2, err := addressStorage.GetOriginalURLByShortURL(short2)
			require.NoError(t, err)
			require.Equal(t, "http://localhost:8080/", fullAddress2)
		} else {
			// require.Equal() ??? TODO
			errlog.Error("Error in charset", zap.Error(err))
		}
	}

}

func TestAddOriginalURLEmptyString(t *testing.T) {
	addressStorage, err := New(false, "")
	if err != nil {
		log.Fatal(err)
	}

	myErr := &storage.EmptyOriginalURLError{}
	_, err = addressStorage.AddOriginalURL("")
	require.Equal(t, err, myErr)
}

func TestGetOriginalURLByShortURL(t *testing.T) {
	addressStorage, err := New(false, "")
	if err != nil {
		log.Fatal(err)
	}

	fullAddress := "http://localhost:8080/"
	name, err := addressStorage.AddOriginalURL(fullAddress)
	require.NoError(t, err)

	short, err := addressStorage.GetOriginalURLByShortURL(name)
	require.NoError(t, err)
	require.Equal(t, fullAddress, short)
}

func TestGetOriginalURLByShortURLUnknownShortURL(t *testing.T) {
	addressStorage, err := New(false, "")
	if err != nil {
		log.Fatal(err)
	}

	unknownShortURL := "aaa"
	_, err = addressStorage.GetOriginalURLByShortURL(unknownShortURL)
	require.Equal(t, err, &storage.NoEntryError{ShortURL: unknownShortURL})
}

// func TestStringWithCharset(t *testing.T) {
// 	newStr, err := helpers.StringWithCharset()
// 	require.NoError(t, err)
// 	require.Len(t, newStr, 10)
// }

// func TestStringWithCharsetInvalidLength(t *testing.T) {
// 	newStr, err := helpers.StringWithCharset() // -10, charSet
// 	require.Equal(t, err, &storage.InvalidLengthError{})
// 	require.Empty(t, newStr)
// }

// func TestStringWithCharsetInvalidCharSet(t *testing.T) {
// 	// charSet := ""

// 	newStr, err := helpers.StringWithCharset() // 10, charSet
// 	require.Equal(t, err, &storage.InvalidCharSetError{})
// 	require.Empty(t, newStr)
// }
