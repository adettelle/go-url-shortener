package urlstorage

import (
	"context"
	"log"

	"github.com/adettelle/go-url-shortener/internal/helpers"
	"github.com/adettelle/go-url-shortener/internal/storage"
	"github.com/adettelle/go-url-shortener/internal/storage/dbstorage"
)

// type CustomerStorage struct {
// 	UserLogin    string
// 	UserPassword string
// 	Adresses     *AddressStorage
// }

type AddressStorage struct {
	// aaa@mail.ru: {aaa: google.com, bbb: ya.ru}
	// Addresses map[string]map[string]string
	Addresses map[string]string // short_url: original_url
	FileName  string            // чтобы можно было синхронно писать изменения в файл FileStoragePath
}

func New(shouldRestore bool, fileStoragePath string) (*AddressStorage, error) {
	if shouldRestore {
		addressStorage, err := ReadJSONFromFile(fileStoragePath) // addressStorage
		if err != nil {
			return nil, err
		}
		addressStorage.FileName = fileStoragePath
		return addressStorage, nil
	}

	addresses := make(map[string]string)

	storage := &AddressStorage{
		Addresses: addresses,
		FileName:  fileStoragePath,
	}
	return storage, nil
}

// возращает полный url по ключу (короткому url)
func (a *AddressStorage) GetOriginalURLByShortURL(shortURL string) (string, error) {
	if addr, ok := a.Addresses[shortURL]; ok {
		return addr, nil
	}

	return "", &storage.NoEntryError{
		ShortURL: shortURL,
	}
}

func (a *AddressStorage) AddOriginalURL(originalURL string) (string, error) {
	if originalURL == "" {
		return "", &storage.EmptyOriginalURLError{}
	}

	randString, err := helpers.StringWithCharset()
	if err != nil {
		return "", err
	}

	a.Addresses[randString] = originalURL

	if a.FileName != "" {
		err := WriteAddressStorageToJSONFile(a.FileName, a)
		if err != nil {
			return "", err
		}
	}

	return randString, nil
}

// отрабатывает завершение приложения (при штатном завершении работы)
// процесс финализации: объекты могут делать работу, пользоваться ресурсамии,
// и при заверщении работы (без работы с БД или с файлом),
// надо содержимое AddressStorage записать на диск (в файл)
func (a *AddressStorage) Finalize() error {
	log.Println("ms.FileName:", a.FileName)
	return WriteAddressStorageToJSONFile(a.FileName, a)
}

func (a *AddressStorage) GetShortURLByOriginalURL(originalURL string) (string, error) {
	for short, origin := range a.Addresses {
		if origin == originalURL {
			return short, storage.NewOriginalURLExistsErr(short, origin)
		}
	}

	return "", nil
}

func (a *AddressStorage) GetAllURLS(ctx context.Context) ([]dbstorage.URL, error) {
	return nil, nil
}
