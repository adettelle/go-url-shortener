package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	defaultAddress    = "localhost:8888"
	defaultURLAddress = "localhost:8000"
)

type Config struct {
	Address    string // отвечает за адрес запуска HTTP-сервера, например, localhost:8888
	URLAddress string // базовый адрес результирующего сокращённого URL
	// (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg)

}

func initFlags() *Config {
	flagAddr := flag.String("a", "", "Net address localhost:port")
	flagURLAddr := flag.String("b", "", "Result url address http://localhost:port/qsd54gFg")

	flag.Parse()

	cfg := Config{
		Address:    getAddr(flagAddr),
		URLAddress: getURLAddr(flagURLAddr),
	}
	return &cfg
}

// сначала проверяем флаги и заполняем структуру конфига оттуда
func New() (*Config, error) {
	cfg := initFlags()

	if cfg.Address == "" {
		cfg.Address = defaultAddress
	}
	if cfg.URLAddress == "" {
		cfg.URLAddress = defaultURLAddress
	}

	ensureAddrFLagIsCorrect(cfg.Address)
	ensureAddrFLagIsCorrect(cfg.URLAddress)

	return cfg, nil
}

func getAddr(flagAddr *string) string {
	log.Println("flagAddress:", *flagAddr, os.Getenv("ADDRESS"))

	addr := os.Getenv("ADDRESS")
	if addr != "" {
		return addr
	}
	return *flagAddr
}

func getURLAddr(flagURLAddr *string) string {
	log.Println("flagURLAddress:", *flagURLAddr, os.Getenv("URL_ADDRESS"))

	addr := os.Getenv("URL_ADDRESS")
	if addr != "" {
		return addr
	}
	return *flagURLAddr
}

func ensureAddrFLagIsCorrect(addr string) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal()
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid port: '%s'", port))
	}
}
