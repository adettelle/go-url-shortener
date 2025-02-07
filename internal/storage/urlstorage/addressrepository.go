package urlstorage

import (
	"log"

	"github.com/adettelle/go-url-shortener/internal/helpers"
	"github.com/adettelle/go-url-shortener/internal/storage"
)

type AddressStorage struct {
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
func (a *AddressStorage) GetAddress(shortURL string) (string, error) {
	if addr, ok := a.Addresses[shortURL]; ok {
		return addr, nil
	}

	return "", &storage.NoEntryError{
		ShortURL: shortURL,
	}
}

func (a *AddressStorage) AddAddress(originalURL string) (string, error) {
	if originalURL == "" {
		return "", &storage.EmptyAddressError{}
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
