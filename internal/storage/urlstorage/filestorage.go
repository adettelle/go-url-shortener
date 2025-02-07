package urlstorage

import (
	"encoding/json"
	"log"
	"os"

	"go.uber.org/zap"
)

// запись в файл AddressStorage
func WriteAddressStorageToJSONFile(storagePath string, addrStorage *AddressStorage) error {
	log.Println("storagePath in WriteAddressStorageToJSONFile:", storagePath)
	// открываем файл для записи (для перезаписи)
	file, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(addrStorage)
	if err != nil {
		return err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()

	logger.Info("writing to file", zap.String("storagePath", storagePath))

	log.Printf("writing to file: %s", storagePath)
	_, err = file.Write(data)
	return err
}
