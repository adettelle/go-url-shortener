package helpers

import (
	"encoding/json"
	"io"
	"os"

	"github.com/adettelle/go-url-shortener/internal/storage"
)

func ReadJSONFromFile(storagePath string, result *storage.AddressStorage) error {
	file, err := os.OpenFile(storagePath, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &result)
}
