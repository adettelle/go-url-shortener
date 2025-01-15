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
	defaultAddress    = "localhost:8080"
	defaultURLAddress = "http://localhost:8080"
)

type Config struct {
	Address    string `envconfig:"SERVER_ADDRESS"` // отвечает за адрес запуска HTTP-сервера, например, localhost:8080
	URLAddress string `envconfig:"BASE_URL"`       // базовый адрес результирующего сокращённого URL
	// (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg)

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

	flag.Parse()

	if cfg.Address == "" {
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

	ensureAddrFLagIsCorrect(cfg.Address)
	ensureStringIsURL(cfg.URLAddress)

	return &cfg, nil
}

func ensureAddrFLagIsCorrect(addr string) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("error in ensuring address: %s", addr))
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid port: '%s'", port))
	}
}

func ensureStringIsURL(addr string) {
	_, err := url.ParseRequestURI(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid url: '%s'", addr))
	}
}
