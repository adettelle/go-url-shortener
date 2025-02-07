package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAddress         = "localhost:8080"
	defaultURLAddress      = "http://localhost:8080"
	defaultFileStoragePath = "/tmp/short-url-db.json"
	// defaultDBParams        = "host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable"
)

// type DBConnection struct {
// 	DBHost     string `envconfig:"DATABASE_HOST" default:"localhost"`
// 	DBPort     string `envconfig:"DATABASE_PORT" default:"5433"`
// 	DBUser     string `envconfig:"DATABASE_USER" default:"postgres"`
// 	DBPassword string `envconfig:"DATABASE_PASSWORD" default:"password"`
// 	DBName     string `envconfig:"DATABASE_NAME" default:"postgres"`
// }

type Config struct {
	Address string `envconfig:"SERVER_ADDRESS"` // отвечает за адрес запуска HTTP-сервера, например, localhost:8080

	URLAddress string `envconfig:"BASE_URL"` // базовый адрес результирующего сокращённого URL
	// (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg)
	FileStoragePath string `envconfig:"FILE_STORAGE_PATH"` // полное имя файла,
	// куда сохраняются данные в формате JSON, пустое значение отключает функцию записи на диск
	DBParams string `envconfig:"DATABASE_DSN"`
	Restore  bool   `json:"restore"` // по умолчанию true
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
		// if cfg.DBParams == "" {
		// 	cfg.DBParams = defaultDBParams
		// }
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

// DBConnStr constructs and returns the PostgreSQL database connection string
// func (dbConn *DBConnection) DBConnStr() string {
// 	dbParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		dbConn.DBHost, dbConn.DBPort, dbConn.DBUser, dbConn.DBPassword, dbConn.DBName)

// 	return dbParams
// }

func (config *Config) ShouldRestore() bool {
	if !config.Restore {
		return false
	}

	// если файл по адресу "/tmp/metrics-db.json" есть, но пустой (fi.Size() == 0),
	// далее проверяем, что именно он записан у нас в config.StoragePath
	// и то стираем всё из хранилища addressStorage, т. к. загружать нечего, ведь файл пустой
	fileStoragePath, err := os.Stat(config.FileStoragePath)
	if err != nil && !os.IsNotExist(err) { // если какая-то другая ошибка
		log.Fatal(err)
	}

	// в этом месте мы знаем, что файл существует, и что Restore = true,
	// значит надо убедится в размере файла
	return fileStoragePath.Size() > 0
}
