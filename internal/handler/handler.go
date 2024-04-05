package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"eCar/config"
	"eCar/internal/domains"
	"eCar/internal/shema"
	"encoding/json"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

type Handler struct {
	service domains.Service
	engine  *gin.Engine
	config  config.Config
}

func NewHandler(service domains.Service, cnf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		engine:  router,
		config:  cnf,
	}
	Route(router, h)
	return h
}

func (s *Handler) Start(ctx context.Context) {
	if s.config.TLS && (s.config.PrivateKey == "" || s.config.Certificate == "") {
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(2024),
			Subject: pkix.Name{
				Organization: []string{"andogeek"},
				Country:      []string{"USA"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Fatal(err)
		}

		certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
		if err != nil {
			log.Fatal(err)
		}

		var certPEM bytes.Buffer
		pem.Encode(&certPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		})

		var privateKeyPEM bytes.Buffer
		pem.Encode(&privateKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
		cerPemFile, err := os.Create("cerPem.crt")
		if err != nil {
			log.Fatalf("Failed to create file: %s", err.Error())
		}
		_, err = cerPemFile.Write(certPEM.Bytes())
		if err != nil {
			log.Fatalf("Failed to write file: %s", err.Error())
		}

		defer cerPemFile.Close()

		privateFile, err := os.Create("private.key")
		if err != nil {
			log.Fatalf("Failed to create file: %s", err.Error())
		}

		_, err = privateFile.Write(privateKeyPEM.Bytes())
		if err != nil {
			log.Fatalf("Failed to write file: %s", err.Error())
		}

		defer privateFile.Close()

		server := &http.Server{
			Addr:    s.config.Host,
			Handler: s.engine.Handler(),
		}

		go func() {
			<-ctx.Done()
			if err := server.Close(); err != nil {
				log.Fatalf("Server forced to shutdown: %v", err)
			}
		}()

		log.Fatal(server.ListenAndServeTLS("cerPem.crt", "private.key"))
	} else if s.config.TLS && s.config.Certificate != "" && s.config.PrivateKey != "" {
		server := &http.Server{
			Addr:    s.config.Host,
			Handler: s.engine.Handler(),
		}

		go func() {
			<-ctx.Done()
			if err := server.Close(); err != nil {
				log.Fatalf("Server forced to shutdown: %v", err)
			}
		}()

		log.Fatal(server.ListenAndServeTLS(s.config.Certificate, s.config.PrivateKey))
	} else {
		server := &http.Server{
			Addr:    s.config.Host,
			Handler: s.engine.Handler(),
		}

		go func() {
			<-ctx.Done()
			if err := server.Close(); err != nil {
				log.Fatalf("Server forced to shutdown: %v", err)
			}
		}()

		log.Fatal(server.ListenAndServe())
	}
}

func (s *Handler) GetNewCars(c *gin.Context) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Здесь мы можем указать ответ, который должен вернуть мок сервер
		rw.Write([]byte(`{"regNum": "X123XX150", "mark": "Lada", "model": "Vesta", "year": 2002, "owner": {"name": "John", "surname": "Doe", "patronymic": "Smith"}}`))
	}))
	// Закрываем сервер, когда он больше не нужен
	defer server.Close()
	var regNums shema.RegNum
	if err := c.ShouldBindJSON(&regNums); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cars []shema.Car

	for _, regNum := range regNums.RegNums {
		resp, err := http.Get(server.URL + "/info?regNum=" + regNum)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			HandlerErr(c, err)
			return
		}

		var car shema.Car
		err = json.Unmarshal(body, &car)
		if err != nil {
			HandlerErr(c, err)
			return
		}

		cars = append(cars, car)
	}

	ctx := c.Request.Context()
	err := s.service.AddCar(ctx, cars)
	if err != nil {
		HandlerErr(c, err)
		return
	}

}
