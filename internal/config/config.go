package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAddress         = "localhost:8080"
	defaultURLAddress      = "http://localhost:8080"
	defaultFileStoragePath = "/tmp/short-url-db.json"
	defaultDBParams        = "host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable"
)

type Config struct {
	Address    string `envconfig:"SERVER_ADDRESS"` // отвечает за адрес запуска HTTP-сервера, например, localhost:8080
	URLAddress string `envconfig:"BASE_URL"`       // базовый адрес результирующего сокращённого URL
	// (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg)
	FileStoragePath string `envconfig:"FILE_STORAGE_PATH"` // полное имя файла,
	// куда сохраняются данные в формате JSON, пустое значение отключает функцию записи на диск
	DBParams string `envconfig:"DATABASE_DSN"`
}

// приоритет:
// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Println("error in process:", err)
		return nil, err
	}

	flagAddr := flag.String("a", "", "Net address localhost:port")
	flagURLAddr := flag.String("b", "", "Result url address http://localhost:port/qsd54gFg")
	flagFileStoragePath := flag.String("f", "", "full file name for data in json format")
	flagDBParams := flag.String("d", "", "db connection params")

	flag.Parse()

	if cfg.Address == "" { // то есть переменная окружения не задана
		cfg.Address = *flagAddr
		if cfg.Address == "" {
			cfg.Address = defaultAddress
		}
	}
	if cfg.URLAddress == "" {
		cfg.URLAddress = *flagURLAddr
		if cfg.URLAddress == "" {
			cfg.URLAddress = defaultURLAddress
		}
	}

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *flagFileStoragePath
		if cfg.FileStoragePath == "" {
			cfg.FileStoragePath = defaultFileStoragePath
		}
	}

	if cfg.DBParams == "" {
		cfg.DBParams = *flagDBParams
		if cfg.DBParams == "" {
			cfg.DBParams = defaultDBParams
		}
	}

	mustBeCorrectAddressFlag(cfg.Address)
	mustBeCorrectURL(cfg.URLAddress)
	// TODO надо ли проверить путь FileStoragePath

	return &cfg, nil
}

func mustBeCorrectAddressFlag(addr string) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("error in ensuring address: %s", addr))
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid port: '%s'", port))
	}
}

func mustBeCorrectURL(addr string) {
	_, err := url.ParseRequestURI(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid url: '%s'", addr))
	}
}
