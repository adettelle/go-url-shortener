package urlstorage

import (
	"encoding/json"
	"io"
	"os"
)

func ReadJSONFromFile(storagePath string) (*AddressStorage, error) {
	// открываем файл для чтения
	jsonFile, err := os.OpenFile(storagePath, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// we initialize our AddressStorage array
	var addrStorage AddressStorage

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'allMetrics' which we defined above
	err = json.Unmarshal(data, &addrStorage)
	if err != nil {
		return nil, err
	}

	return &addrStorage, nil
}
