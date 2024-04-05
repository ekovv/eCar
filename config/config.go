package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host        string `json:"host"`
	TLS         bool   `json:"tls"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private"`
	DB          string `json:"dsn"`
}

type F struct {
	host        *string
	tls         *bool
	certificate *string
	privateKey  *string
	db          *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=")
	f.tls = flag.Bool("tls", false, "-tls=")
	f.certificate = flag.String("cert", "", "-cert=path_to_certificate")
	f.privateKey = flag.String("key", "", "-key=path_to_key")
	f.db = flag.String("d", "", "-d=db")

}

func New() (c Config) {
	flag.Parse()

	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if envHost := os.Getenv("HOST"); envHost != "" {
		f.host = &envHost
	}
	if envCert := os.Getenv("CERTIFICATE"); envCert != "" {
		f.certificate = &envCert
	}
	if envKey := os.Getenv("PRIVATE_KEY"); envKey != "" {
		f.privateKey = &envKey
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	if envTLS := os.Getenv("TLS"); envTLS != "" {
		tlsValue, err := strconv.ParseBool(envTLS)
		if err != nil {
			log.Fatalf("Error parsing TLS value: %v", err)
		}
		f.tls = &tlsValue
	}
	c.Host = *f.host
	c.TLS = *f.tls
	c.Certificate = *f.certificate
	c.PrivateKey = *f.privateKey
	c.DB = *f.db
	return c

}
