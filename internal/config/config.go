package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

const (
	defaultAddress    = "localhost:8080"
	defaultURLAddress = "http://localhost:8080"
)

type Config struct {
	Address    string `envconfig:"SERVER_ADDRESS"` //  default:"localhost:8080" отвечает за адрес запуска HTTP-сервера, например, localhost:8080
	URLAddress string `envconfig:"BASE_URL"`       //  default:"http://localhost:8080" базовый адрес результирующего сокращённого URL
	// (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg)

}

// func initFlags() *Config {
// 	flagAddr := flag.String("a", "", "Net address localhost:port")
// 	flagURLAddr := flag.String("b", "", "Result url address http://localhost:port/qsd54gFg")

// 	flag.Parse()

// 	cfg := Config{
// 		Address:    getAddr(flagAddr),
// 		URLAddress: getURLAddr(flagURLAddr),
// 	}
// 	return &cfg
// }

// приоритет:
// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
func New() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil { // TODO будет ли перезапись!!!!
		log.Println("error in process:", err)
		return nil, err
	}

	flagAddr := flag.String("a", "", "Net address localhost:port")
	flagURLAddr := flag.String("b", "", "Result url address http://localhost:port/qsd54gFg")

	flag.Parse() // TODO будет ли перезапись!!!!

	if cfg.Address == "" {
		cfg.Address = getAddr(flagAddr)
		if cfg.Address == "" {
			cfg.Address = defaultAddress
		}
	}
	if cfg.URLAddress == "" {
		cfg.URLAddress = getURLAddr(flagURLAddr)
		if cfg.URLAddress == "" {
			cfg.URLAddress = defaultURLAddress
		}
	}

	ensureAddrFLagIsCorrect(cfg.Address)
	// ensureAddrFLagIsCorrect(cfg.URLAddress)

	return &cfg, nil
}

func getAddr(flagAddr *string) string {
	// log.Println("flagAddress:", *flagAddr, os.Getenv("ADDRESS"))

	addr := os.Getenv("ADDRESS")
	if addr != "" {
		return addr
	}
	return *flagAddr
}

func getURLAddr(flagURLAddr *string) string {
	// log.Println("flagURLAddress:", *flagURLAddr, os.Getenv("URL_ADDRESS"))

	addr := os.Getenv("URL_ADDRESS")
	if addr != "" {
		return addr
	}
	return *flagURLAddr
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
